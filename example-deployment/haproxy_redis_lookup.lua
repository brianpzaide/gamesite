local redis = require("redis")

function get_backend(txn)
    local redis_host = "db"
    local redis_port = 6379

    local client = redis.connect(redis_host, redis_port)
    
    local path = txn.sf:path()
    print(path)

    local room_id = path:match("^/gamesite/rooms/([%w%-]+)/ws$")
    if not room_id then
      room_id = path:match("^/gamesite/rooms/([%w%-]+)$")
    end

    if not room_id then
        print("no roomid")
        txn:set_var("txn.backend", "default_backend")
        return
    end

    local backend_server = client:get(room_id)
    print(backend_server)

    if backend_server then
        txn:set_var("txn.backend", backend_server)
    else
        txn:set_var("txn.backend", "default_backend")
    end
end

core.register_action("get_backend", { "http-req" }, get_backend, 0)