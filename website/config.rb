set :base_url, "https://www.ottoproject.io/"

activate :hashicorp do |h|
  h.name        = "otto"
  h.version     = "0.1.2"
  h.github_slug = "hashicorp/otto"

  h.minify_javascript = false
end
