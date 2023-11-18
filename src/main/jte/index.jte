@import gg.jte.Content

@import io.github.jakedegiovanni.comicshelf.views.Index

@param Index page
@param Content content
@param Content footer

<html lang="en">
<head>
    <meta charset="UTF-8">

    <title>Comicshelf</title>

    <link href="/static/index.css" rel="stylesheet"  />
    <script src="/static/index.js"></script>

    <link href='https://fonts.googleapis.com/css?family=Roboto' rel='stylesheet'>

    <script src="https://unpkg.com/htmx.org@1.9.4" integrity="sha384-zUfuhFKKZCbHTY6aRR46gxiqszMk5tcHjsVFxnUo8VMus4kHGVdIYVbOYYNlKmHV" crossorigin="anonymous"></script>
</head>
<body onload="activateNavItems()">
<div class="bar navbar">
    <form method="get">
        <label for="date">Release Week</label>
        <input type="date" id="date" name="date" value="${page.getDate()}" />
        <button type="submit">Submit</button>
    </form>

    <div>Marvel Unlimited ComicShelf</div>

    <div class="navigation">
        <a class="nav-item" id="/marvel-unlimited/comics" href="/marvel-unlimited/comics?date=${page.getDate()}">Comics</a>
        <a class="nav-item" id="/marvel-unlimited/series" href="/marvel-unlimited/series?date=${page.getDate()}">Series</a>
    </div>
</div>

<div class="main">
    ${content}
</div>

<div class="pusher"></div>

<div class="bar bottombar">
    ${footer}
</div>
</body>
</html>