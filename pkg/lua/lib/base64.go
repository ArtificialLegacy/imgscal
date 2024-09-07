package lib

import (
	"encoding/base64"

	"github.com/ArtificialLegacy/imgscal/pkg/log"
	"github.com/ArtificialLegacy/imgscal/pkg/lua"
	golua "github.com/yuin/gopher-lua"
)

const LIB_BASE64 = "base64"

/// @lib Base64
/// @import base64
/// @desc
/// Utility library for encoding and decoding base64 data.

func RegisterBase64(r *lua.Runner, lg *log.Logger) {
	lib, tab := lua.NewLib(LIB_BASE64, r, r.State, lg)

	/// @func encode(data, url?, raw?) -> string
	/// @arg data {string}
	/// @arg? url {bool} - If true, use URL encoding.
	/// @arg? raw {bool} - If true, use raw encoding.
	/// @returns {string} - The base64 encoded data.
	lib.CreateFunction(tab, "encode",
		[]lua.Arg{
			{Type: lua.STRING, Name: "data"},
			{Type: lua.BOOL, Name: "url", Optional: true},
			{Type: lua.BOOL, Name: "raw", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			data := args["data"].(string)
			raw := args["raw"].(bool)

			var out string
			if !args["url"].(bool) {
				if raw {
					out = base64.RawStdEncoding.EncodeToString([]byte(data))
				} else {
					out = base64.StdEncoding.EncodeToString([]byte(data))
				}
			} else {
				if raw {
					out = base64.RawURLEncoding.EncodeToString([]byte(data))
				} else {
					out = base64.URLEncoding.EncodeToString([]byte(data))
				}
			}

			state.Push(golua.LString(out))
			return 1
		})

	/// @func decode(data, url?, raw?) -> string
	/// @arg data {string}
	// @arg? url {bool} - If true, use URL encoding.
	/// @arg? raw {bool} - If true, use raw encoding.
	/// @returns {string} - The base64 decoded data.
	lib.CreateFunction(tab, "decode",
		[]lua.Arg{
			{Type: lua.STRING, Name: "data"},
			{Type: lua.BOOL, Name: "url", Optional: true},
			{Type: lua.BOOL, Name: "raw", Optional: true},
		},
		func(state *golua.LState, d lua.TaskData, args map[string]any) int {
			data := args["data"].(string)
			raw := args["raw"].(bool)

			var out []byte
			var err error
			if !args["url"].(bool) {
				if raw {
					out, err = base64.RawStdEncoding.DecodeString(data)
				} else {
					out, err = base64.StdEncoding.DecodeString(data)
				}
			} else {
				if raw {
					out, err = base64.RawURLEncoding.DecodeString(data)
				} else {
					out, err = base64.URLEncoding.DecodeString(data)
				}
			}

			if err != nil {
				lua.Error(state, lg.Appendf("failed to decode base64 data: %s", log.LEVEL_ERROR, err))
			}

			state.Push(golua.LString(out))
			return 1
		})
}
