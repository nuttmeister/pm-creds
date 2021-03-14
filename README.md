# pm-creds

pm-creds is a middle-ware between Postman and your credentials provider and securely sets your credentials as environmental variables on your requests or collections.

## How It Works

More info coming ...

## Providers

Currently the following providers are fully supported.

|Provider|Method|Comment|
|-|-|-|
|AWS|redentials+config files or default evaluation chain|Supports permanent and temporary profiles stored in the credentials file. As well as fetching using default provider chain with profile name $default.|

## Installation

You can either download any of the pre-built binaries from the releases page or build it yourself.

### Downloading pre-built

You will find the pre-build versions for `macOS`, `linux` and `windows` under [releases](/releases).

### Building

#### Requirements

- `git`
- `make`
- `go`

#### Clone and Build

Run the following commands to clone, build and install.

```shell
git clone git@github.com:nuttmeister/pm-creds.git
cd pm-creds
make build
sudo make install
pm-creds --cert
```

### Running

To run `pm-creds` you will first have needed to generate certificates. You can do that by following the instructions on [Generate Certificates](#generate-certificates).

#### Generate Certificates

First time you use `pm-creds` you will need to generate `CA Certificate`, `Server Certificate` and a `Client Certificate`.

```shell
pm-creds --cert
```

This will generate certificates valid for *10 years*.  
If you for some reason need to use this program with `Safari` and not `Postman` you can specify a custom validity time specified in days.

```shell
pm-creds --cert --validity 365
```

If you need to recreate the certificates for some reason due to an error or that they have expired you will either need to run the command with the overwrite flag. Since otherwise the program will exit in error if the files already exist.

```shell
pm-creds --cert --overwrite
```


### Postman

You will need to configure Postman to use `pm-creds` properly by installing it's `CA Certificate` as well as an `Client Certificate`.

#### Custom certificate

Go to `Settings -> Certificates` and add the `cert/minica.pem` to `CA Certificates`.

Then add a `Client Certificate` with the following settings:

`Host`: https://`localhost`:`9999`  
`CRT file`: `cert/localhost/cert.pem`  
`KEY file`: `cert/localhost/key.pem`  

#### Adding profile to the environment

Either create a new environment in Postman or edit a current one and add the following variable.
This will control what aws profile you will use for any request made with this environment active.

`aws_profile`: `<the aws profile from ./aws/credentials you want to use>`

#### Adding AWS Auth to Request / Collection

Then choose the `AWS Signature` under `Authorization` on either the `Collection` or on the single `Request` with
the following settings.

`AccessKey`: `{{aws_access_key_id}}`  
`SecretKey`: `{{aws_secret_access_key}}`  
`AWS Region`: `<set if you need it>`  
`Service Name`: `<set if you need it>`  
`Session Token`: `{{aws_session_token}}`  

#### Add Pre Request script to Request / Collection

Then create the following `pre-request script` on either the `Collection` or on the single `Request`.

```js
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

```

## Usage

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

### Start proxy

Running `pm-creds` with no options will load config and providers form the standard config directory `~/.pm-creds/`.

### Run Postman Request

Run an `Request` that is configured with the Auth and Pre-request script on it or on the collection
as described above.

Once you hit send go to the console window and either authorize the request or not.
