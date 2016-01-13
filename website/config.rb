set :base_url, "https://www.ottoproject.io/"

activate :hashicorp do |h|
  h.name        = "otto"
  h.version     = "0.2.0"
  h.github_slug = "hashicorp/otto"
end
