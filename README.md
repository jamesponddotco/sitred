# SitRed

**SitRed** is a simple, yet powerful, server that fetches a
website's sitemap, parses it to extract the URLs, and then redirects
incoming HTTP requests to a random URL from the list.

The server is written in Go and was designed for sitemaps generated by
the [Yoast SEO WordPress
plugin](https://wordpress.org/plugins/wordpress-seo/), but as long as
you don't feed it a sitemap index, you should be able to use with it
other kinds of sitemaps. Additionally, it should work efficiently even
with large sitemaps.

## Usage

Hosting your own instance of **SitRed** is easy:

* [Hosting the service](doc/hosting.md)

Using **SitRed** is even easier:

* [Using the service](doc/using.md)

## Installation

### From source

First install the dependencies:

- Go 1.21 or above.
- make.
- [scdoc](https://git.sr.ht/~sircmpwn/scdoc).

Then compile and install:

```bash
make
sudo make install
```

## Contributing

Anyone can help make **SitRed** better. Send patches on the [mailing
list](https://lists.sr.ht/~jamesponddotco/sitred-devel) and report bugs
on the [issue tracker](https://todo.sr.ht/~jamesponddotco/sitred).

You must sign-off your work using `git commit --signoff`. Follow the
[Linux kernel developer's certificate of
origin](https://www.kernel.org/doc/html/latest/process/submitting-patches.html#sign-your-work-the-developer-s-certificate-of-origin)
for more details.

All contributions are made under [the EUPL license](LICENSE.md).

## Resources

The following resources are available:

- [Support and general discussions](https://lists.sr.ht/~jamesponddotco/sitred-discuss).
- [Patches and development related questions](https://lists.sr.ht/~jamesponddotco/sitred-devel).
- [Instructions on how to prepare patches](https://git-send-email.io/).
- [Feature requests and bug reports](https://todo.sr.ht/~jamesponddotco/sitred).

---

Released under the [EUPL License](LICENSE.md).