#-------------------------------------------------------------------------
# Configure Middleman
#-------------------------------------------------------------------------

helpers do
  def livestream_active?
    # Must set key for date
    ENV["LIVESTREAM_ACTIVE"].present?
  end
end

set :base_url, "https://www.ottoproject.io/"

activate :hashicorp do |h|
  h.version         = ENV["OTTO_VERSION"]
  h.bintray_enabled = ENV["BINTRAY_ENABLED"] == "1"
  h.bintray_repo    = "mitchellh/otto"
  h.bintray_user    = "mitchellh"
  h.bintray_key     = ENV["BINTRAY_API_KEY"]
  h.github_slug     = "hashicorp/otto"

  h.minify_javascript = false
end
