package rbac

default allow = false

# Admin can do anything
allow {
    input.role == "admin"
}

# User can access /api/public
allow {
    input.role == "user"
    input.path == "/api/public"
    input.method == "GET"
}

# User cannot access /api/admin (implicit by default allow = false)
