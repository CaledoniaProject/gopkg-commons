package controllers

import (
	"net/http"
	"strings"

	"github.com/CaledoniaProject/gopkg-commons/commons"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
)

var (
	sessionKeyRedirectURL = "redirectURL"
	sessionKeyOAuthState  = "oauthState"
	sessionKeyUserInfo    = "userInfo"
)

func OAuthLoginHandler(w http.ResponseWriter, r *http.Request, appConfig *commons.StandardConfig) {
	var (
		state       = commons.RandomString(10)
		vars        = mux.Vars(r)
		configBlock = commons.GetOAuthConfigBlock(appConfig.OAuthConfigList, vars["provider"])
	)

	// 校验oauth类型
	if configBlock == nil {
		http.Error(w, "invalid provider", http.StatusBadRequest)
		return
	}

	// 选择性增加重定向地址
	if tmp := r.URL.Query().Get("redirectURL"); strings.HasPrefix(tmp, "/") {
		appConfig.SessionManager.Put(r.Context(), sessionKeyRedirectURL, tmp)
	}

	// 存储state
	appConfig.SessionManager.Put(r.Context(), sessionKeyOAuthState, state)
	http.Redirect(w, r, configBlock.GetConfig().AuthCodeURL(state), http.StatusFound)
}

func OAuthCallbackHandler(w http.ResponseWriter, r *http.Request, standardConfig *commons.StandardConfig, contextLogger *logrus.Logger, convertUserInfo func(*commons.OAuthUserInfo) (any, error)) {
	var (
		code        = r.URL.Query().Get("code")
		state       = r.URL.Query().Get("state")
		vars        = mux.Vars(r)
		provider    = vars["provider"]
		configBlock = commons.GetOAuthConfigBlock(standardConfig.OAuthConfigList, provider)
	)

	// 校验oauth类型
	if configBlock == nil {
		http.Error(w, "invalid provider", http.StatusBadRequest)
		return
	}

	// 防止CSRF
	if state != standardConfig.SessionManager.PopString(r.Context(), sessionKeyOAuthState) {
		http.Error(w, "state mismatch", http.StatusBadRequest)
		return
	}

	// 获取用户信息
	oauthUser, err := configBlock.GetUserInfo(r.Context(), provider, code)
	if err != nil {
		contextLogger.Errorf("get user info: %v", err)
		http.Error(w, "Internal error. Please try again later.", http.StatusBadRequest)
		return
	} else if oauthUser.Email == "" {
		// github可能邮箱为空
		// contextLogger.Errorf("empty email: %v", oauthUser)
		// http.Error(w, "Internal error. Please try again later.", http.StatusBadRequest)
		// return
	}

	// 处理用户登录
	if userInfo, err := convertUserInfo(oauthUser); err != nil {
		http.Error(w, "Internal error. Please try again later.", http.StatusBadRequest)
		return
	} else {
		// 创建或者更新用户
		if tx := standardConfig.SqlDB.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(userInfo); tx.Error != nil {
			contextLogger.Errorf("create user: %v", err)
			http.Error(w, "Internal error. Please try again later.", http.StatusBadRequest)
			return
		}

		if resetter, ok := userInfo.(commons.PrimaryKeyResetter); ok {
			resetter.ResetPrimaryKey()
		}

		// 查询用户信息: 有的接口不提供邮箱，所以根据ID；防止ID已经存在，所以用另外一个空的model
		if tx := standardConfig.SqlDB.Model(userInfo).
			Where("Provider = ? AND ExternalId = ?", oauthUser.Provider, oauthUser.Id).
			Take(userInfo); tx.Error != nil {
			contextLogger.Errorf("find user by %s and %s: %v", oauthUser.Provider, oauthUser.Id, tx.Error)
			http.Error(w, "Internal error. Please try again later.", http.StatusBadRequest)
			return
		}

		// 创建session
		standardConfig.SessionManager.Put(r.Context(), sessionKeyUserInfo, userInfo)
	}

	// 重定向
	redirectURL := standardConfig.SessionManager.PopString(r.Context(), sessionKeyRedirectURL)
	if redirectURL == "" {
		redirectURL = "/"
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
}
