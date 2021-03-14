# pm-creds

pm-creds is a middle-ware between Postman and your credentials provider and securely sets your credentials as environmental variables on your requests or collections.


## Providers

Currently the following providers are fully supported.

|Provider|Method|Comment|
|-|-|-|
|AWS|Profiles (credentials and config files)|Supports permanent and temporary profiles stored in the credentials file.|
|AWS|Default evaluation chain|Supports fetching using default provider chain when using profile name $default.|

## Gettings started

We recommend you download one of the pre-built binaries.  
You will find them for `macOS`, `linux` and `windows` under [releases](https://github.com/nuttmeister/pm-creds/releases).

### Building

You can build using go (version 1.16 or above) by running `cd pm-creds && go build`.

### Running

#### Usage help

```text
pm-creds --help
Usage of pm-creds:
  --config-dir string
        Location of the config files (default "/home/user/.pm-creds")
  --create-certs
        If certificates should be generated
  --create-config
        If the default config should be created
  --overwrite
        If new config/certificates should overwrite old
```

#### Generate config and certificates

To get started you will need to generate `default config` and `certificates`.  
To do this run the commands below.

```shell
pm-creds --create-config --create-certs
```

If you config and / or certificates are broken for some reason you can add the flag `--overwrite`
and `--create-config` and/or `--create-certs` will allow you to overwrite the already existing files.

#### Running

To run the proxy just start it with `pm-creds` and wait for it to start listening.  
It's possible to use a custom config directory, then specify the directory with the `--config-dir` option.


### Postman

You will need to configure Postman to use `pm-creds` properly by installing it's `CA Certificate` as well as the `Server Certificate`.

#### Configure Certificates

Go to `Settings -> Certificates`

##### CA Certificate

Add the `certs/ca-cert.pem` as Postmans `CA Certificate`. (default: `~/.pm-creds/certs/ca-cert.pem`).

##### Client Certificate

Then add a `Client Certificate` with the following settings.

```text
Host:     https://localhost:9999
CRT file: certs/server-cert.pem
KEY file: certs/server-key.pem
```

#### Adding profile to the environment

Either create a new environment in Postman or edit a current one and add the following variable.
This will control what aws profile you will use for any request made with this environment active.

```text
aws_profile: <the aws profile you want to use>
```

#### Adding AWS Auth to Request / Collection

Then choose the `AWS Signature` under `Authorization` on either the `Collection` or on the single `Request` with
the following settings.

```text
AccessKey    : {{aws_access_key_id}}
SecretKey    : {{aws_secret_access_key}}
AWS Region   : <set if you need it>
Service Name : <set if you need it>
Session Token: {{aws_session_token}}
```

#### Add Pre Request script to Request / Collection

Then create the following `pre-request script` on either the `Collection` or on the single `Request`.

To be sure to use the latest version of this script and for scripts for other providers please have a look in the `/postman` directory.

**Below example is for AWS**

```js
const profile = pm.environment.get("aws_profile")
if (!profile) {
    throw new Error("'aws_profile' variable not set")
}

pm.sendRequest({
    url: `https://localhost:9999/aws/${profile}`,
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
```

#### Run Postman Request

Run an `Request` that is configured with the Auth and Pre-request script on it or on the collection
as described above.

Once you hit send go to the console window and either authorize the request or not.


## Configuration

You can have a look at the `config.default.toml` file for the default configuration that will be created
when running with the `--create-config` option.

The default config directory is `~/.pm-creds`.
