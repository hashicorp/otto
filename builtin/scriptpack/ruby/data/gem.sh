# ruby_gemfile_contains checks if a Gemfile has a certain gem in it
ruby_gemfile_contains() {
    local name=$1

    if [ -f Gemfile.lock ]; then
        grep -e " $name (" Gemfile.lock > /dev/null
        return $?
    fi

    if [ -f Gemfile ]; then
        grep -e "gem .$name." Gemfile > /dev/null
        return $?
    fi

    return 1
}

# ruby_gemfile_apt installs packages for Gems that are detected.
ruby_gemfile_apt() {
    _ruby_gemfile_queue=()
    _ruby_gemfile_check curb "libcurl3 libcurl3-gnutls libcurl4-openssl-dev"
    _ruby_gemfile_check capybara-webkit "libqt4-dev"
    _ruby_gemfile_check mysql2 "libmysqlclient-dev"
    _ruby_gemfile_check nokogiri "zlib1g-dev"
    _ruby_gemfile_check pg "libpq-dev"
    _ruby_gemfile_check rmagick "libmagickwand-dev"
    _ruby_gemfile_check sqlite3 "libsqlite3-dev"
    _ruby_gemfile_check libxml-ruby "libxml2-dev"
    _ruby_gemfile_check paperclip "imagemagick"
    _ruby_gemfile_check poltergeist "phantomjs"
    _ruby_gemfile_check tiny_tds "freetds-dev"

    if [ -n "${_ruby_gemfile_queue-}" ]; then
        otto_output "Installing native gem system dependencies..."
        apt_update_once
        apt_install "${_ruby_gemfile_queue[@]}"
    fi
}

# Internal functions for accumulating the queue of things to install
# for a Gemfile.
_ruby_gemfile_queue=()
_ruby_gemfile_check() {
    local gem=$1
    local deps=$2

    if ruby_gemfile_contains $gem; then
        otto_output "Detected the gem: ${gem}"
        _ruby_gemfile_queue+=($deps)
    fi
}
