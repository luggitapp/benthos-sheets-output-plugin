# Benthos Sheets Output Plugin

This benthos plugin provides a simple way to output to a google spreadsheet.


## Getting Started
For this plugin to work you will need to get a service account credential from google cloud.

1. Create a service account https://console.cloud.google.com/identity/serviceaccounts/create with the `Project -> Editor` role.
2. Create a key for the service account, download it and save it.
3. Create a new sheet if you don't have it already.
4. Share the sheet with your new service-account email as an editor.
### Running and testing

Clone the project and after doing the steps above you need to define the `GOOGLE_APPLICATION_CREDENTIALS` environment variable in your terminal.

You can do this by entering the command bellow in your terminal

`export GOOGLE_APPLICATION_CREDENTIALS="<path to your service account key>`

Now let's run the example provided in the repository.

1. Install deps `go mod vendor`
2. Replace the sheetId in `examples/stdin.yaml` with your sheet id found in the url.
3. Run `make run` and enter `{"field1": "test1", "field2": "test2"}` in the running process.
4. You should now see one row with `test1,test2`
5. Thats it.

## Plugin Configuration
For now the configuration is pretty simple.

### Parameters

- SheetId: The sheet id of your spreesheet this can be found on the url.
![](https://miro.medium.com/max/1000/1*1QhoxsBaUr65GUqkRibJLA.png)

### Example benthos output configuration
```yaml
  output:
    type: sheets
    plugin:
      sheetId: "1z5W9eXXzr9fG0lYcTtXSgtYfnusIPCuAIomsZvhkO5s"
```


## Deployment
We use helm to deploy our containers with that we provide a chart for deploying your first processor to sheets.

### Usage
1. Clone this repository
2. Edit the `values.yaml` value to your needs
3. Obtain a google service account like said above and create a secret in the same namespace you are deploying
```kubectl create secret generic google-service-account --from-file=sva.json=/path/to/svc-key.json -n processors```
1. Install the chart `helm install myprocessor chart/ -f values.yaml -n processors`
2. Done



## Authors
Hugo Fonseca <hugo.fonseca@luggit.app>


## LICENSE
MIT