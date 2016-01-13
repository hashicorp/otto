# Otto Website

This subdirectory contains the entire source for the
[Otto Website](https://www.ottoproject.io/). This is a
[Middleman](http://middlemanapp.com) project, which builds a static site from
these source files.

## Contributions Welcome!

If you find a typo or you feel like you can improve the HTML, CSS, or
JavaScript, we welcome contributions. Feel free to open issues or pull requests
like any normal GitHub project, and we will merge it in.

## Running the Site Locally

Running the site locally is simple.

1. Install [Otto](https://www.ottoproject.io/download.html)
1. Clone this repo to your local machine
1. Change into the `website` directory
1. Run:

        $ otto compile

1. Run:

        $ otto dev

1. This will output the static IP address of the instance. Now you can SSH into
   the instance and run the static site:

        $ otto dev ssh

1. To run the site as a server:

        $ bundle exec middleman server

1. To compile the static HTML:

        $ bundle exec middleman build
