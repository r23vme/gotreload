## How to use: 

https://github.com/r23vme/gotreload/blob/master/gotreload_test.go

## Echo endpoint
```
gr := gotreload.New()
router.GET(gr.URL, func(c echo.Context) error {
  gr.ServeHTTP(c.Response().Writer, c.Request())
  return nil
})
```

## Echo middlewares:

### Inject `<script>` after `</html>`
Have to be the last Middleware to modify the response to be really "after `</html>`".
```
router := echo.New()

gr := gotreload.New()
echoMiddleware := func(next echo.HandlerFunc) echo.HandlerFunc {
  return func(c echo.Context) error {
    if err := next(c); err != nil {
      return err
    }
    c.Response().Writer.Write([]byte(gr.Script()))
    return nil
  }
}

router.Use(echoMiddleware)
```

### Inject `<script>` in `<body></body>`
```
type interceptorWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (r interceptorWriter) Write(b []byte) (int, error) {
	return r.body.Write(b)
}

func main(){
  router := echo.New()

  gr := gotreload.New()
  echoMiddlewareInBody := func(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
      writer := &interceptorWriter{body: &bytes.Buffer{}, ResponseWriter: c.Response().Writer}
      c.Response().Writer = writer

      if err := next(c); err != nil {
        return err
      }

      body := writer.body.Bytes()
      if b := body; bytes.Contains(b, []byte("</body>")) {
        b = bytes.ReplaceAll(b, []byte("</body>"), []byte(fmt.Sprintf("%s</body>", gr.Script())))
        body = b
      }

      writer.ResponseWriter.Write(body)
      return nil
    }
  }
  router.Use(echoMiddlewareInBody)
}
```
