# GCP Resource Tree Builder

This program allows you to build a tree representation of your Google Cloud Platform (GCP) folders and their associated projects.

## Features

- Build a tree representation of GCP folders and projects
- Collects projects for each folder
- Concurrently processes folders for faster execution

## Requirements

- Go version 1.20 or later
- Google Cloud Platform Service Account JSON key file with appropriate permissions

## Usage

1. Clone this repository and navigate to the project directory.
2. Build the binary by running go build.
3. Execute the binary with the following flags:

```bash
./[binary_name] --key-file [path_to_service_account_key] --folders [comma_separated_folder_IDs]
```

- Replace [binary_name] with the binary name generated after building the project.
- Replace [path_to_service_account_key] with the path to your GCP Service Account JSON key file.
- Replace [comma_separated_folder_IDs] with a comma-separated list of GCP folder IDs you want to build the tree for.

Example:

```bash
./gcp-folder-project-tree-builder --key-file /path/to/keyfile.json --folders 123456789012,234567890123

```

The output will be a JSON representation of the folder tree, including the folder names, IDs, children folders, and projects associated with each folder. The JSON result will be printed to the console.

Example JSON output:

```json
[
  {
    "name": "folder1",
    "id": "123456789012",
    "children": [
      {
        "name": "subfolder1",
        "id": "234567890123",
        "projects": [
          {
            "name": "projects/my-project-1",
            "project_id": "my-project-1",
            "state": "ACTIVE",
            "create_time": "2021-08-25T16:26:31.882Z"
          }
        ]
      }
    ],
    "projects": [
      {
        "name": "projects/my-project-2",
        "project_id": "my-project-2",
        "state": "ACTIVE",
        "create_time": "2023-03-30T16:26:31.882Z"
      }
    ]
  },
  {
    "name": "folder2",
    "id": "234567890123",
    "children": [],
    "projects": [
      {
        "name": "projects/my-project-3",
        "project_id": "my-project-3",
        "state": "ACTIVE",
        "create_time": "2023-03-30T16:26:31.882Z"
      }
    ]
  }
]
```

## Pre-built binaries

You can download pre-built binaries for Windows, Linux, and macOS from the [releases page](https://github.com/david1155/gcp-folder-project-tree-builder/releases).

## Troubleshooting

If you encounter any issues while running the program, please ensure that:

1. You have provided the correct path to the GCP Service Account JSON key file.
2. The Service Account has sufficient permissions to access the resource manager API for the specified folders and projects.
3. You have provided a valid comma-separated list of folder IDs.
4. You have the required version of Go installed to build and run the program.

If you still experience issues, consider checking the logs for error messages or unexpected behavior. The logs may provide more context about the problem and help you to identify the root cause.

Feel free to raise an issue on the repository if you need further assistance.