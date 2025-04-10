---@param workflow imgscal_WorkflowInit
function init(workflow)
	workflow.import({
		"net",
		"cli",
	})
end

function main()
	net.route("/hello", function(w, r)
		if r.method ~= net.METHOD_GET then
			w.write_header(net.STATUS_METHOD_NOT_ALLOWED)
		else
			w.write("world\n")
		end
	end)

	local addr = ":3131"

	cli.println("listening on " .. addr)
	_ = net.serve(addr)
end
