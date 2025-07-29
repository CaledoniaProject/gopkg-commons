package markdown

type MarkdownHeader struct {
	Title       string   `yaml:"title"`
	Excerpt     string   `yaml:"excerpt"`
	AuthorURL   []string `yaml:"author_url"`
	PublishDate string   `yaml:"publish_date"`
	Tags        []string `yaml:"tags"`
}
