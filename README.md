# keyrun
Cli wrapper managing env keys and encrypting files

#### Why

Cli tools sometimes are not flexible or secure enough managing env variables and local files. Keyrun is a tiny cli wrapper, that decrypts files and sets up env variables, executes the command then cleans the things up. It is desgned for terraform but can be reused for any cli tool.

#### Getting started
1. Place config .keyrun.yml in the directory with terraform recipes:
```
env:
  AWS_ACCESS_KEY: aws_access_key
  AWS_SECRET_KEY: aws_secret_key
key: tf_key
```
aws_access_key/aws_secret_key/tf_key can be any name, that's how it will be stored in your OS keychain, prefixed with `keyrun_`.

2. Create all required keys in keychain:
```
keyrun key create
```

3. Encrypt state files:
```
keyrun encrypt terraform.tfstate terraform.tfstate.enc
keyrun encrypt terraform.tfstate.backup terraform.tfstate.backup.enc
```

4. No 4. Now you can run any command you need as usual:
```
keyrun exec -- terraform plan
```
