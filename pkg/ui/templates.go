package ui

var getRunbookTplString = `<div class="row">
  <div class="col-xl">
    <div class="runbook">
{{{RunbookHTML}}}
    </div>
  </div>
</div>
<script>
$(".runbook h1,h2,h3,h4,h5,h6").each((index, val) => {
  const $el = $(val);
  const id = $el.attr("id");
  const text = $el.text();
  $el.text("");
  $el.append('<a href="#' + id + '">' + text + '</a>');
});
</script>
`

var layoutTplString = `<html>
  <head>
    <title>
      Runbooks
    </title>
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.5.0/css/bootstrap.min.css" integrity="sha384-9aIt2nRpC12Uk9gS9baDl411NQApFmC26EwAOH8WgZl5MYYxFfc+NcPb1dKGj7Sk" crossorigin="anonymous">
    <link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/highlight.js/10.1.1/styles/default.min.css">
    <script src="//cdnjs.cloudflare.com/ajax/libs/highlight.js/10.1.1/highlight.min.js"></script>
    <script>hljs.initHighlightingOnLoad();</script>
    <script src="//code.jquery.com/jquery-3.5.1.slim.min.js"></script>
  </head>
  <body>
    <header>
      <nav class="navbar navbar-expand-lg navbar-dark bg-dark mb-4">
        <a class="navbar-brand" href="/runbooks">Runbooks</a>
      </nav>
    </header>
    <div class="container">
      {{{content}}}
    </div>
  </body>
</html>
`

var listRunbooksTplString = `<div class="row">
  <div class="col-xl">
    <table class="table">
      <thead>
        <tr>
          <th>
            Name
          </th>
        </tr>
      </thead>
      <tbody>
  {{#Runbooks}}
        <tr>
          <td>
            <a href="/runbooks/{{Name}}">{{Name}}</a>
          </td>
        </tr>
  {{/Runbooks}}
      </tbody>
    </table>
  </div>
</div>
`
