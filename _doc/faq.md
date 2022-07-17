[HTML5 Boilerplate homepage](https://html5boilerplate.com/) | [Documentation
table of contents](TOC.md)

# Frequently asked questions

* [Why is the Google Analytics code at the bottom? Google recommends it be
  placed in the
  `<head>`.](#why-is-the-google-analytics-code-at-the-bottom-google-recommends-it-be-placed-in-the-head)


---

## Why is the Google Analytics code at the bottom? Google recommends it be placed in the `<head>`.

The main advantage of placing it in the `<head>` is that you will track the
user's `pageview` even if they leave the page before it has been fully loaded.

Here's a handy quote from [Mathias
Bynens](https://mathiasbynens.be/notes/async-analytics-snippet#comment-50) about
our placement choice.
>I should point out that it’s Google — not me — recommending to place this
script before all other scripts in the document. The only real advantage is to
catch a pageView call if your page fails to load completely (for example, if the
user aborts loading, or quickly closes the page, etc.). Personally, I wouldn’t
count that as a page view, so I actually prefer to place this script at the
bottom, after all other scripts. This keeps all the scripts together and
reinforces that scripts at the bottom are the right move. (Usually I concatenate
and minify all my scripts into one .js file — the GA snippet being the suffix.)

