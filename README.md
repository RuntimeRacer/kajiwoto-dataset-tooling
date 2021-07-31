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
  - [ ] Milestone 2: diff (point out and store differences) and sync (eliminate differences) of CSV files and datasets.
  
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
go build -o kajitool
```
That's it. The commands should execute without any issued. Now you've got your `kajitool` binary ready. 


## License & Copyright notice
- `kajitool` is free software licensed under the [Apache-2.0 License](LICENSE).
- [Kajiwoto](https://kajiwoto.com/) is a platform for creating AI companions. The author of `kajitool` is in no way aligned with Kajiwoto or paid for his work. The sole purpose of kajitool is to support his own AI development efforts and the Kajiwoto community. 
