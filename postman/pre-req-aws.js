const profile = pm.environment.get("aws_profile")
if (!profile) {
    throw new Error("'aws_profile' variable not set")
}

pm.sendRequest({
    url: `https://localhost:9999/${profile}`,
    method: "POST",
    }, function (_, response) {
        if (response.status == "OK") {
            const body = response.json()
            pm.variables.set("aws_access_key_id", body.accessKey)
            pm.variables.set("aws_secret_access_key", body.secretKey)
            if (body.sessionToken) {
                pm.variables.set("aws_session_token", body.sessionToken)
            }
            console.log(`using aws credentials from '${profile}'`)
            return
        } else {
            throw new Error(response.text() || "unknown error fetching aws credentials")
        }
    }
)
