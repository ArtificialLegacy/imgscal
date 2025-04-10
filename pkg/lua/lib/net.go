package lib

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/yuin/gopher-lua"
)

const LIB_NET = "net"

/// @lib Networking
/// @import net
/// @desc
/// Library for making http requests and serving endpoints.

func RegisterNet(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_NET, r, r.State, lg)

	/// @func req(method, url, body, header, contentLength?) -> struct<net.Response>
	/// @arg method {string<net.Method>}
	/// @arg url {string}
	/// @arg body {string}
	/// @arg header {struct<net.Values>}
	/// @arg? contentLength {int}
	/// @returns {struct<net.Response>}
	lib.CreateFunction(tab, "req",
		[]lua.Arg{
			{Type: lua.STRING, Name: "method"},
			{Type: lua.STRING, Name: "url"},
			{Type: lua.STRING, Name: "body"},
			{Type: lua.RAW_TABLE, Name: "header"},
			{Type: lua.INT, Name: "contentLength", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			method := args["method"].(string)
			urlstr := args["url"].(string)
			body := args["body"].(string)
			header := netvaluesBuild(args["header"].(*golua.LTable))
			contentLength := args["contentLength"].(int)

			req, err := http.NewRequest(method, urlstr, strings.NewReader(body))
			if err != nil {
				lua.Error(state, lg.Appendf("failed to create request: %s", log.LEVEL_ERROR, err))
			}
			req.Header = header
			if contentLength >= 0 {
				req.ContentLength = int64(contentLength)
			} else {
				req.ContentLength = -1
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				lua.Error(state, lg.Appendf("failed to make http request to url %s: with error (%s)", log.LEVEL_WARN, urlstr, err))
			}
			defer resp.Body.Close()

			t := responseTable(lib, state, lg, resp)

			state.Push(t)
			return 1
		})

	/// @func req_get(url) -> struct<net.Response>
	/// @arg url {string}
	/// @returns {struct<net.Response>}
	lib.CreateFunction(tab, "req_get",
		[]lua.Arg{
			{Type: lua.STRING, Name: "url"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			urlstr := args["url"].(string)
			resp, err := http.Get(urlstr)
			if err != nil {
				lua.Error(state, lg.Appendf("failed to make http GET request to url %s: with error (%s)", log.LEVEL_WARN, urlstr, err))
			}
			defer resp.Body.Close()

			t := responseTable(lib, state, lg, resp)

			state.Push(t)
			return 1
		})

	/// @func req_head(url) -> struct<net.Response>
	/// @arg url {string}
	/// @returns {struct<net.Response>}
	lib.CreateFunction(tab, "req_head",
		[]lua.Arg{
			{Type: lua.STRING, Name: "url"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			urlstr := args["url"].(string)
			resp, err := http.Head(urlstr)
			if err != nil {
				lua.Error(state, lg.Appendf("failed to make http HEAD request to url %s: with error (%s)", log.LEVEL_WARN, urlstr, err))
			}
			defer resp.Body.Close()

			t := responseTable(lib, state, lg, resp)

			state.Push(t)
			return 1
		})

	/// @func req_post(url, contentType, body) -> struct<net.Response>
	/// @arg url {string}
	/// @arg contentType {string}
	/// @arg body {string}
	/// @returns {struct<net.Response>}
	lib.CreateFunction(tab, "req_post",
		[]lua.Arg{
			{Type: lua.STRING, Name: "url"},
			{Type: lua.STRING, Name: "contentType"},
			{Type: lua.STRING, Name: "body"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			urlstr := args["url"].(string)
			contentType := args["contentType"].(string)
			body := args["body"].(string)

			resp, err := http.Post(urlstr, contentType, strings.NewReader(body))
			if err != nil {
				lua.Error(state, lg.Appendf("failed to make http POST request to url %s: with error (%s)", log.LEVEL_WARN, urlstr, err))
			}
			defer resp.Body.Close()

			t := responseTable(lib, state, lg, resp)

			state.Push(t)
			return 1
		})

	/// @func req_post_form(url, values) -> struct<net.Response>
	/// @arg url {string}
	/// @arg values {struct<net.Values>}
	/// @returns {struct<net.Response>}
	lib.CreateFunction(tab, "req_post_form",
		[]lua.Arg{
			{Type: lua.STRING, Name: "url"},
			{Type: lua.RAW_TABLE, Name: "values"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			urlstr := args["url"].(string)
			values := netvaluesBuild(args["values"].(*golua.LTable))

			resp, err := http.PostForm(urlstr, url.Values(values))
			if err != nil {
				lua.Error(state, lg.Appendf("failed to make http POST request to url %s: with error (%s)", log.LEVEL_WARN, urlstr, err))
			}
			defer resp.Body.Close()

			t := responseTable(lib, state, lg, resp)

			state.Push(t)
			return 1
		})

	/// @func values() -> struct<net.Values>
	/// @returns {struct<net.Values>}
	lib.CreateFunction(tab, "values",
		[]lua.Arg{},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			return 1
		})

	/// @func canonical_header_key(key) -> string
	/// @arg key {string}
	/// @returns {string}
	lib.CreateFunction(tab, "canonical_header_key",
		[]lua.Arg{
			{Type: lua.STRING, Name: "key"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			canon := http.CanonicalHeaderKey(args["key"].(string))

			state.Push(golua.LString(canon))
			return 1
		})

	/// @func detect_content(content) -> string
	/// @arg content {string}
	/// @returns {string}
	/// @desc
	/// Detects the http content of a data represented as a string.
	lib.CreateFunction(tab, "detect_content",
		[]lua.Arg{
			{Type: lua.STRING, Name: "content"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			contentType := http.DetectContentType([]byte(args["content"].(string)))

			state.Push(golua.LString(contentType))
			return 1
		})

	/// @func status_text(code) -> string
	/// @arg code {int<net.Status>}
	/// @returns {string}
	lib.CreateFunction(tab, "status_text",
		[]lua.Arg{
			{Type: lua.INT, Name: "code"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			text := http.StatusText(args["code"].(int))

			state.Push(golua.LString(text))
			return 1
		})

	/// @func route(pattern, fn)
	/// @arg pattern {string}
	/// @arg fn {function(w struct<net.ResponseWriter>, r struct<net.Request>)}
	lib.CreateFunction(tab, "route",
		[]lua.Arg{
			{Type: lua.STRING, Name: "pattern"},
			{Type: lua.FUNC, Name: "fn"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			pattern := args["pattern"].(string)
			fn := args["fn"].(*golua.LFunction)

			routeThread, _ := state.NewThread()

			http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
				respWriter := responseWriterTable(lib, routeThread, lg, w)
				request := requestTable(lib, routeThread, lg, r)

				routeThread.Push(fn)
				routeThread.Push(respWriter)
				routeThread.Push(request)
				routeThread.Call(2, 0)
			})

			return 0
		})

	/// @func serve(addr)
	/// @arg addr {string}
	/// @blocking
	lib.CreateFunction(tab, "serve",
		[]lua.Arg{
			{Type: lua.STRING, Name: "addr"},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			server := http.Server{Addr: args["addr"].(string)}
			closed := make(chan struct{})

			go func() {
				defer func() {
					if p := recover(); p != nil {
						closed <- struct{}{}
					}
				}()

				err := server.ListenAndServe()
				if err != nil {
					if !errors.Is(err, http.ErrServerClosed) {
						lua.Error(state, lg.Appendf("failed to start server: %s", log.LEVEL_ERROR, err))
					}
				}

				closed <- struct{}{}
			}()

			select {
			case <-r.Ctx.Done():
				_ = server.Close()
			case <-closed:
			}

			return 0
		})

	/// @constants Method {string}
	/// @const METHOD_GET
	/// @const METHOD_HEAD
	/// @const METHOD_POST
	/// @const METHOD_PUT
	/// @const METHOD_PATCH
	/// @const METHOD_DELETE
	/// @const METHOD_CONNECT
	/// @const METHOD_OPTIONS
	/// @const METHOD_TRACE
	tab.RawSetString("METHOD_GET", golua.LString(http.MethodGet))
	tab.RawSetString("METHOD_HEAD", golua.LString(http.MethodHead))
	tab.RawSetString("METHOD_POST", golua.LString(http.MethodPost))
	tab.RawSetString("METHOD_PUT", golua.LString(http.MethodPut))
	tab.RawSetString("METHOD_PATCH", golua.LString(http.MethodPatch))
	tab.RawSetString("METHOD_DELETE", golua.LString(http.MethodDelete))
	tab.RawSetString("METHOD_CONNECT", golua.LString(http.MethodConnect))
	tab.RawSetString("METHOD_OPTIONS", golua.LString(http.MethodOptions))
	tab.RawSetString("METHOD_TRACE", golua.LString(http.MethodTrace))

	/// @constants Status {int}
	/// @const STATUS_CONTINUE
	/// @const STATUS_SWITCHING_PROTOCOLS
	/// @const STATUS_PROCESSING
	/// @const STATUS_EARLY_HINTS
	/// @const STATUS_OK
	/// @const STATUS_CREATED
	/// @const STATUS_ACCEPTED
	/// @const STATUS_NON_AUTHORITATIVE_INFO
	/// @const STATUS_NO_CONTENT
	/// @const STATUS_RESET_CONTENT
	/// @const STATUS_PARTIAL_CONTENT
	/// @const STATUS_MULTI_STATUS
	/// @const STATUS_ALREADY_REPORTED
	/// @const STATUS_IM_USED
	/// @const STATUS_MULTIPLE_CHOICES
	/// @const STATUS_MOVED_PERMANENTLY
	/// @const STATUS_FOUND
	/// @const STATUS_SEE_OTHER
	/// @const STATUS_NOT_MODIFIED
	/// @const STATUS_USE_PROXY
	/// @const STATUS_TEMPORARY_REDIRECT
	/// @const STATUS_PERMANENT_REDIRECT
	/// @const STATUS_BAD_REQUEST
	/// @const STATUS_UNAUTHORIZED
	/// @const STATUS_PAYMENT_REQUIRED
	/// @const STATUS_FORBIDDEN
	/// @const STATUS_NOT_FOUND
	/// @const STATUS_METHOD_NOT_ALLOWED
	/// @const STATUS_NOT_ACCEPTABLE
	/// @const STATUS_PROXY_AUTH_REQUIRED
	/// @const STATUS_REQUEST_TIMEOUT
	/// @const STATUS_CONFLICT
	/// @const STATUS_GONE
	/// @const STATUS_LENGTH_REQUIRED
	/// @const STATUS_PRECONDITION_FAILED
	/// @const STATUS_PAYLOAD_TOO_LARGE
	/// @const STATUS_URI_TOO_LONG
	/// @const STATUS_UNSUPPORTED_MEDIA_TYPE
	/// @const STATUS_RANGE_NOT_SATISFIABLE
	/// @const STATUS_EXPECTATION_FAILED
	/// @const STATUS_TEAPOT
	/// @const STATUS_MISDIRECTED_REQUEST
	/// @const STATUS_UNPROCESSABLE_ENTITY
	/// @const STATUS_LOCKED
	/// @const STATUS_FAILED_DEPENDENCY
	/// @const STATUS_TOO_EARLY
	/// @const STATUS_UPGRADE_REQUIRED
	/// @const STATUS_PRECONDITION_REQUIRED
	/// @const STATUS_TOO_MANY_REQUESTS
	/// @const STATUS_REQUEST_HEADER_FIELDS_TOO_LARGE
	/// @const STATUS_UNAVAILABLE_FOR_LEGAL_REASONS
	/// @const STATUS_INTERNAL_SERVER_ERROR
	/// @const STATUS_NOT_IMPLEMENTED
	/// @const STATUS_BAD_GATEWAY
	/// @const STATUS_SERVICE_UNAVAILABLE
	/// @const STATUS_GATEWAY_TIMEOUT
	/// @const STATUS_HTTP_VERSION_NOT_SUPPORTED
	/// @const STATUS_VARIANT_ALSO_NEGOTIATES
	/// @const STATUS_INSUFFICIENT_STORAGE
	/// @const STATUS_LOOP_DETECTED
	/// @const STATUS_NOT_EXTENDED
	/// @const STATUS_NETWORK_AUTH_REQUIRED
	tab.RawSetString("STATUS_CONTINUE", golua.LNumber(http.StatusContinue))
	tab.RawSetString("STATUS_SWITCHING_PROTOCOLS", golua.LNumber(http.StatusSwitchingProtocols))
	tab.RawSetString("STATUS_PROCESSING", golua.LNumber(http.StatusProcessing))
	tab.RawSetString("STATUS_EARLY_HINTS", golua.LNumber(http.StatusEarlyHints))
	tab.RawSetString("STATUS_OK", golua.LNumber(http.StatusOK))
	tab.RawSetString("STATUS_CREATED", golua.LNumber(http.StatusCreated))
	tab.RawSetString("STATUS_ACCEPTED", golua.LNumber(http.StatusAccepted))
	tab.RawSetString("STATUS_NON_AUTHORITATIVE_INFO", golua.LNumber(http.StatusNonAuthoritativeInfo))
	tab.RawSetString("STATUS_NO_CONTENT", golua.LNumber(http.StatusNoContent))
	tab.RawSetString("STATUS_RESET_CONTENT", golua.LNumber(http.StatusResetContent))
	tab.RawSetString("STATUS_PARTIAL_CONTENT", golua.LNumber(http.StatusPartialContent))
	tab.RawSetString("STATUS_MULTI_STATUS", golua.LNumber(http.StatusMultiStatus))
	tab.RawSetString("STATUS_ALREADY_REPORTED", golua.LNumber(http.StatusAlreadyReported))
	tab.RawSetString("STATUS_IM_USED", golua.LNumber(http.StatusIMUsed))
	tab.RawSetString("STATUS_MULTIPLE_CHOICES", golua.LNumber(http.StatusMultipleChoices))
	tab.RawSetString("STATUS_MOVED_PERMANENTLY", golua.LNumber(http.StatusMovedPermanently))
	tab.RawSetString("STATUS_FOUND", golua.LNumber(http.StatusFound))
	tab.RawSetString("STATUS_SEE_OTHER", golua.LNumber(http.StatusSeeOther))
	tab.RawSetString("STATUS_NOT_MODIFIED", golua.LNumber(http.StatusNotModified))
	tab.RawSetString("STATUS_USE_PROXY", golua.LNumber(http.StatusUseProxy))
	tab.RawSetString("STATUS_TEMPORARY_REDIRECT", golua.LNumber(http.StatusTemporaryRedirect))
	tab.RawSetString("STATUS_PERMANENT_REDIRECT", golua.LNumber(http.StatusPermanentRedirect))
	tab.RawSetString("STATUS_BAD_REQUEST", golua.LNumber(http.StatusBadRequest))
	tab.RawSetString("STATUS_UNAUTHORIZED", golua.LNumber(http.StatusUnauthorized))
	tab.RawSetString("STATUS_PAYMENT_REQUIRED", golua.LNumber(http.StatusPaymentRequired))
	tab.RawSetString("STATUS_FORBIDDEN", golua.LNumber(http.StatusForbidden))
	tab.RawSetString("STATUS_NOT_FOUND", golua.LNumber(http.StatusNotFound))
	tab.RawSetString("STATUS_METHOD_NOT_ALLOWED", golua.LNumber(http.StatusMethodNotAllowed))
	tab.RawSetString("STATUS_NOT_ACCEPTABLE", golua.LNumber(http.StatusNotAcceptable))
	tab.RawSetString("STATUS_PROXY_AUTH_REQUIRED", golua.LNumber(http.StatusProxyAuthRequired))
	tab.RawSetString("STATUS_REQUEST_TIMEOUT", golua.LNumber(http.StatusRequestTimeout))
	tab.RawSetString("STATUS_CONFLICT", golua.LNumber(http.StatusConflict))
	tab.RawSetString("STATUS_GONE", golua.LNumber(http.StatusGone))
	tab.RawSetString("STATUS_LENGTH_REQUIRED", golua.LNumber(http.StatusLengthRequired))
	tab.RawSetString("STATUS_PRECONDITION_FAILED", golua.LNumber(http.StatusPreconditionFailed))
	tab.RawSetString("STATUS_PAYLOAD_TOO_LARGE", golua.LNumber(http.StatusRequestEntityTooLarge))
	tab.RawSetString("STATUS_URI_TOO_LONG", golua.LNumber(http.StatusRequestURITooLong))
	tab.RawSetString("STATUS_UNSUPPORTED_MEDIA_TYPE", golua.LNumber(http.StatusUnsupportedMediaType))
	tab.RawSetString("STATUS_RANGE_NOT_SATISFIABLE", golua.LNumber(http.StatusRequestedRangeNotSatisfiable))
	tab.RawSetString("STATUS_EXPECTATION_FAILED", golua.LNumber(http.StatusExpectationFailed))
	tab.RawSetString("STATUS_TEAPOT", golua.LNumber(http.StatusTeapot))
	tab.RawSetString("STATUS_MISDIRECTED_REQUEST", golua.LNumber(http.StatusMisdirectedRequest))
	tab.RawSetString("STATUS_UNPROCESSABLE_ENTITY", golua.LNumber(http.StatusUnprocessableEntity))
	tab.RawSetString("STATUS_LOCKED", golua.LNumber(http.StatusLocked))
	tab.RawSetString("STATUS_FAILED_DEPENDENCY", golua.LNumber(http.StatusFailedDependency))
	tab.RawSetString("STATUS_TOO_EARLY", golua.LNumber(http.StatusTooEarly))
	tab.RawSetString("STATUS_UPGRADE_REQUIRED", golua.LNumber(http.StatusUpgradeRequired))
	tab.RawSetString("STATUS_PRECONDITION_REQUIRED", golua.LNumber(http.StatusPreconditionRequired))
	tab.RawSetString("STATUS_TOO_MANY_REQUESTS", golua.LNumber(http.StatusTooManyRequests))
	tab.RawSetString("STATUS_REQUEST_HEADER_FIELDS_TOO_LARGE", golua.LNumber(http.StatusRequestHeaderFieldsTooLarge))
	tab.RawSetString("STATUS_UNAVAILABLE_FOR_LEGAL_REASONS", golua.LNumber(http.StatusUnavailableForLegalReasons))
	tab.RawSetString("STATUS_INTERNAL_SERVER_ERROR", golua.LNumber(http.StatusInternalServerError))
	tab.RawSetString("STATUS_NOT_IMPLEMENTED", golua.LNumber(http.StatusNotImplemented))
	tab.RawSetString("STATUS_BAD_GATEWAY", golua.LNumber(http.StatusBadGateway))
	tab.RawSetString("STATUS_SERVICE_UNAVAILABLE", golua.LNumber(http.StatusServiceUnavailable))
	tab.RawSetString("STATUS_GATEWAY_TIMEOUT", golua.LNumber(http.StatusGatewayTimeout))
	tab.RawSetString("STATUS_HTTP_VERSION_NOT_SUPPORTED", golua.LNumber(http.StatusHTTPVersionNotSupported))
	tab.RawSetString("STATUS_VARIANT_ALSO_NEGOTIATES", golua.LNumber(http.StatusVariantAlsoNegotiates))
	tab.RawSetString("STATUS_INSUFFICIENT_STORAGE", golua.LNumber(http.StatusInsufficientStorage))
	tab.RawSetString("STATUS_LOOP_DETECTED", golua.LNumber(http.StatusLoopDetected))
	tab.RawSetString("STATUS_NOT_EXTENDED", golua.LNumber(http.StatusNotExtended))
	tab.RawSetString("STATUS_NETWORK_AUTH_REQUIRED", golua.LNumber(http.StatusNetworkAuthenticationRequired))
}

func netvaluesTable(lib *lua.Lib, state *golua.LState, values map[string][]string) *golua.LTable {
	/// @struct Values
	/// @method add(self, key string, value string) -> self
	/// @method del(self, key string) -> self
	/// @method get(key string) -> string
	/// @method get_all(key string) -> []string
	/// @method has(key string) -> bool
	/// @method set(self, key string, value string) -> self
	/// @method set_multi(self, key string, values []string) -> self

	t := state.NewTable()

	valuesTable := state.NewTable()
	for k, v := range values {
		vlist := state.NewTable()
		for i, s := range v {
			vlist.RawSetInt(i+1, golua.LString(s))
		}

		valuesTable.RawSetString(k, vlist)
	}
	t.RawSetString("__values", valuesTable)

	lib.BuilderFunction(state, t, "add",
		[]lua.Arg{
			{Type: lua.STRING, Name: "key"},
			{Type: lua.STRING, Name: "value"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			key := args["key"].(string)
			value := args["value"].(string)

			vtable := t.RawGetString("__values").(*golua.LTable)
			vvalue := vtable.RawGetString(key)

			if vlist, ok := vvalue.(*golua.LTable); ok {
				vlist.Append(golua.LString(value))
				vvalue = vlist
			} else {
				vlist := state.NewTable()
				vlist.Append(golua.LString(value))
				vvalue = vlist
			}

			vtable.RawSetString(key, vvalue)
		})

	lib.BuilderFunction(state, t, "del",
		[]lua.Arg{
			{Type: lua.STRING, Name: "key"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			key := args["key"].(string)

			vtable := t.RawGetString("__values").(*golua.LTable)
			vtable.RawSetString(key, golua.LNil)
		})

	lib.TableFunction(state, t, "get",
		[]lua.Arg{
			{Type: lua.STRING, Name: "key"},
		},
		func(state *golua.LState, args map[string]any) int {
			key := args["key"].(string)
			var result string

			vtable := t.RawGetString("__values").(*golua.LTable)
			vvalue := vtable.RawGetString(key)
			if vlist, ok := vvalue.(*golua.LTable); ok {
				if vlist.Len() > 0 {
					result = string(vlist.RawGetInt(1).(golua.LString))
				}
			}

			state.Push(golua.LString(result))
			return 1
		})

	lib.TableFunction(state, t, "get_all",
		[]lua.Arg{
			{Type: lua.STRING, Name: "key"},
		},
		func(state *golua.LState, args map[string]any) int {
			key := args["key"].(string)
			result := state.NewTable()

			vtable := t.RawGetString("__values").(*golua.LTable)
			vvalue := vtable.RawGetString(key)
			if vlist, ok := vvalue.(*golua.LTable); ok {
				for i := range vlist.Len() {
					result.RawSetInt(i+1, vlist.RawGetInt(i+1))
				}
			}

			state.Push(result)
			return 1
		})

	lib.TableFunction(state, t, "has",
		[]lua.Arg{
			{Type: lua.STRING, Name: "key"},
		},
		func(state *golua.LState, args map[string]any) int {
			key := args["key"].(string)

			vtable := t.RawGetString("__values").(*golua.LTable)
			_, ok := vtable.RawGetString(key).(*golua.LTable)

			state.Push(golua.LBool(ok))
			return 1
		})

	lib.BuilderFunction(state, t, "set",
		[]lua.Arg{
			{Type: lua.STRING, Name: "key"},
			{Type: lua.STRING, Name: "value"},
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			key := args["key"].(string)
			value := args["value"].(string)

			vlist := state.NewTable()
			vlist.Append(golua.LString(value))

			vtable := t.RawGetString("__values").(*golua.LTable)
			vtable.RawSetString(key, vlist)
		})

	lib.BuilderFunction(state, t, "set_multi",
		[]lua.Arg{
			{Type: lua.STRING, Name: "key"},
			lua.ArgArray("values", lua.ArrayType{Type: lua.STRING}, false),
		},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			key := args["key"].(string)
			value := args["values"].([]any)

			vlist := state.NewTable()
			for i, v := range value {
				vlist.RawSetInt(i+1, golua.LString(v.(string)))
			}

			vtable := t.RawGetString("__values").(*golua.LTable)
			vtable.RawSetString(key, vlist)
		})

	return t
}

func netvaluesBuild(t *golua.LTable) map[string][]string {
	values := map[string][]string{}

	vvalue := t.RawGetString("__values")
	if vlist, ok := vvalue.(*golua.LTable); ok {
		vlist.ForEach(func(k, v golua.LValue) {
			vtable := v.(*golua.LTable)
			varray := make([]string, vtable.Len())

			for i := range vtable.Len() {
				varray[i] = string(vtable.RawGetInt(i + 1).(golua.LString))
			}

			values[string(k.(golua.LString))] = varray
		})
	}

	return values
}

func responseTable(lib *lua.Lib, state *golua.LState, lg *log.Logger, resp *http.Response) *golua.LTable {
	/// @struct Response
	/// @prop body {string}
	/// @prop status {string}
	/// @prop status_code {int<net.Status>}
	/// @prop header {struct<net.Values>}
	/// @prop content_length {int}
	/// @prop uncompressed {bool}
	/// @prop proto {string}
	/// @prop proto_major {int}
	/// @prop proto_minor {int}

	t := state.NewTable()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		lua.Error(state, lg.Appendf("failed to read body: %s", log.LEVEL_ERROR, err))
	}
	t.RawSetString("body", golua.LString(body))

	t.RawSetString("status", golua.LString(resp.Status))
	t.RawSetString("status_code", golua.LNumber(resp.StatusCode))

	t.RawSetString("header", netvaluesTable(lib, state, resp.Header))

	t.RawSetString("content_length", golua.LNumber(resp.ContentLength))
	t.RawSetString("uncompressed", golua.LBool(resp.Uncompressed))

	t.RawSetString("proto", golua.LString(resp.Proto))
	t.RawSetString("proto_major", golua.LNumber(resp.ProtoMajor))
	t.RawSetString("proto_minor", golua.LNumber(resp.ProtoMinor))

	return t
}

func responseWriterTable(lib *lua.Lib, state *golua.LState, lg *log.Logger, w http.ResponseWriter) *golua.LTable {
	/// @struct ResponseWriter
	/// @method header() -> struct<net.Values>
	/// @method write(data string) -> int
	/// @method write_header(code int<net.Status>)

	t := state.NewTable()

	lib.TableFunction(state, t, "header",
		[]lua.Arg{},
		func(state *golua.LState, args map[string]any) int {
			header := w.Header()

			state.Push(netvaluesTable(lib, state, header))
			return 1
		})

	lib.TableFunction(state, t, "write",
		[]lua.Arg{
			{Type: lua.STRING, Name: "data"},
		},
		func(state *golua.LState, args map[string]any) int {
			data := args["data"].(string)
			n, err := w.Write([]byte(data))
			if err != nil {
				lua.Error(state, lg.Appendf("failed to write data to response: %s", log.LEVEL_ERROR, err))
			}

			state.Push(golua.LNumber(n))
			return 1
		})

	lib.TableFunction(state, t, "write_header",
		[]lua.Arg{
			{Type: lua.INT, Name: "code"},
		},
		func(state *golua.LState, args map[string]any) int {
			code := args["code"].(int)
			w.WriteHeader(code)

			return 0
		})

	return t
}

func requestTable(lib *lua.Lib, state *golua.LState, lg *log.Logger, r *http.Request) *golua.LTable {
	/// @struct Request
	/// @prop method {string<net.Method>}
	/// @prop proto {string}
	/// @prop proto_major {int}
	/// @prop proto_minor {int}
	/// @prop header {struct<net.Values>}
	/// @prop body {string}
	/// @prop content_length {int}
	/// @prop host {string}
	/// @prop remote_addr {string}
	/// @prop request_uri {string}
	/// @prop form {struct<net.Values>}
	/// @prop form_post {struct<net.Values>}
	/// @method parse_form(self) -> self

	t := state.NewTable()

	t.RawSetString("method", golua.LString(r.Method))

	t.RawSetString("proto", golua.LString(r.Proto))
	t.RawSetString("proto_major", golua.LNumber(r.ProtoMajor))
	t.RawSetString("proto_minor", golua.LNumber(r.ProtoMinor))

	header := netvaluesTable(lib, state, r.Header)
	t.RawSetString("header", header)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		lua.Error(state, lg.Appendf("failed to read body: %s", log.LEVEL_ERROR, err))
	}
	t.RawSetString("body", golua.LString(body))

	t.RawSetString("content_length", golua.LNumber(r.ContentLength))
	t.RawSetString("host", golua.LString(r.Host))
	t.RawSetString("remote_addr", golua.LString(r.RemoteAddr))
	t.RawSetString("request_uri", golua.LString(r.RequestURI))

	t.RawSetString("form", golua.LNil)
	t.RawSetString("form_post", golua.LNil)

	lib.BuilderFunction(state, t, "parse_form",
		[]lua.Arg{},
		func(state *golua.LState, t *golua.LTable, args map[string]any) {
			err := r.ParseForm()
			if err != nil {
				lua.Error(state, lg.Appendf("failed to parse form data: %s", log.LEVEL_ERROR, err))
			}

			form := netvaluesTable(lib, state, r.Form)
			t.RawSetString("form", form)
			formPost := netvaluesTable(lib, state, r.PostForm)
			t.RawSetString("form_post", formPost)
		})

	return t
}
