# kajiwoto-dataset-tooling
![kajitool logo](doc/gopher-kajiwoto.png)

## News
#### 2021-07-31
MVP Milestone 1 reached! It is now possible to use `kajitool` for simple uploading and downloading of dataset content via CSV files.

## General Info

`kajitool` is a small CLI application to enable dataset developers the ability to gain more control over the contents of their datasets than the existing tooling provided by Kajiwoto currently offers. 

Goals of this project:
- Enabling the ability to download Kajiwoto Datasets and adding training data to existing datasets.
- Enabling the ability to sync datasets and/or local files with training data.
- Enabling the ability to point out differences between datasets and/or local files.

Non-Goals:
- Remote management of Dataset meta information (e.g. rename, delete, change settings etc.)
- 

### Roadmap
- Minimum viable product (MVP)
  - [x] Milestone 1 (Achieved on 2021-07-31): Simple download and upload of datasets
  - [ ] Milestone 2: 
    - diff (point out and store differences) and sync (eliminate differences) of CSV files and datasets.
    - Support for context information of training data.
  
- Main Release (v1.0)
  - [ ] Support for additional export options (e.g. JSON, SQlite)
  

## USAGE

### Prequisites
- [Git Client](https://git-scm.com/)
- [Golang 1.16 or higher](https://golang.org/dl/)

### Building `kajitool`
Building `kajitool` is easy. First, open a console window. Then execute following commands.
```
git clone https://github.com/RuntimeRacer/kajiwoto-dataset-tooling
cd kajiwoto-dataset-tooling

# NIX-Users
go build -o kajitool

# WIN-Users
go build -o kajitool.exe
```
That's it. The commands should execute without any issued. Now you've got your `kajitool` binary ready.

### Login to Kajiwoto using `kajitool`
Next step is to perform the login. This has to be done once after building the binary, and each time your server session expires. However, Kajitool will store your credentials and session key (after successful login) locally at `$HOME\.kajitool.yaml`. This way you'd only have to issue `./kajitool login` to renew your session once expired.
```
# NIX-Users
./kajitool login -u '$USERNAME' -p '$PASSWORD'
# WIN-Users
kajitool.exe login -u '$USERNAME' -p '$PASSWORD'
```
Once logged in, you can download any dataset of your own, free ones, or the ones that you've purchased on the marketplace. Currently, only storing them in `.csv` files is supported. For further info on how the data has to be read, please check the detailed explaination in the comment of type [DatasetEntry](/blob/main/cmd/dataset.go#L65).

### Downloading a dataset using `kajitool`
The `download` command currently only supports downloading if an exact dataset ID is provided, and the output format will always be a `.csv` file. To retrieve the dataset ID, navigate to your Dataset via web app. The URL should be something like `https://kajiwoto.com/d/XXX`, where `XXX` is the ID of your dataset. Copy this value and provide it as the source param for the `download` command.
```
# NIX-Users
./kajitool dataset download -s '$DATASET_ID' -t 'dataset.csv'
# WIN-Users
kajitool.exe dataset download -s '$DATASET_ID' -t 'dataset.csv'
```
Depending on the size of the dataset, `kajitool` might have to issue multiple requests to fetch all entries. The progress will be printed in the console window. To not hammer Kajiwoto's API too much, there is a small pause between those requests.

Another feature of the `download` command, is that it will check for duplicate entries in the dataset on the fly. If it encounters a duplicate, `kajitool` will print a warning and also store the duplicate IDs along with the dataset entries in the resulting `.csv` file.

### Uploading training data to a dataset using `kajitool`
The `upload` command currently only supports uploading if an exact dataset ID is provided, and the source format will always be expected to be a `.csv` file. To retrieve the dataset ID, navigate to your Dataset via web app. The URL should be something like `https://kajiwoto.com/d/XXX`, where `XXX` is the ID of your dataset. Copy this value and provide it as the target param for the `upload` command.
```
# NIX-Users
./kajitool dataset upload -t '$DATASET_ID' -s 'dataset.csv'
# WIN-Users
kajitool.exe dataset upload -t '$DATASET_ID' -s 'dataset.csv'
```
When uploading, `kajitool` will analyze the provided source data and will only attempt uploading data from lines which have no ID provided. This is a simple safety mechanism to make sure that it's not uploading duplicates into an existing dataset, without having to download that dataset first. So if you're attempting to feed a new Dataset from existing ones that you downloaded previously, make sure to empty the ID Row (first one) of the source file.

Right now, despite I believe the API supports batch uploading of training data (as the `multi` param of the graphQL Mutation indicates), `upload` issues the training requests one by one. To not hammer Kajiwoto's API too much, there is a small pause between those requests.

Also, the upload functionality currently is somewhat limited. Despite the source data may contain chat history context information, it is currently not supported to add the context for trainings when uploading. The reason for this is simply, that there might be multiple follow-up dialog options for an initial starter dialog (Example: You let your Kaji Waifu know that you want to eat something; she asks you what you want to eat; and next dialog responses then depend on your answer and kaji mood etc.)


## License & Copyright notice
- `kajitool` is free software licensed under the [Apache-2.0 License](LICENSE).
- [Kajiwoto](https://kajiwoto.com/) is a platform for creating AI companions. The author of `kajitool` is in no way aligned with Kajiwoto or paid for his work. The sole purpose of kajitool is to support his own AI development efforts and the Kajiwoto community. 
