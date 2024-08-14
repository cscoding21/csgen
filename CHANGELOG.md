## v0.6.7 (2024-08-14)

### Fix

- update CLA handler replacement with struct package

## v0.6.6 (2024-08-14)

### Fix

- update default package when invalid one exists

## v0.6.5 (2024-08-13)

### Fix

- fix handling of command-line-arguments

## v0.6.4 (2024-08-12)

### Fix

- added nil handling for module name inference

## v0.6.3 (2024-08-12)

### Fix

- updated loadmodule function to accept parameters

## v0.6.2 (2024-08-08)

### Refactor

- update method for returning a package or struct from the associated receivers

## v0.6.1 (2024-08-08)

### Fix

- add ID to module in package

## v0.6.0 (2024-08-07)

### Feat

- add loadmodule function to get a comprehensive graph of a full go module

## v0.5.0 (2024-06-30)

### Feat

- add new file header generator for one-off generated files and a function to infer a package name from a path

## v0.4.3 (2024-06-25)

### Fix

- fix the GetFile method to property handle absolute paths

## v0.4.2 (2024-06-13)

### Fix

- re-update GetFileName and NewCSGenBuilderForFile to conform to go standards

## v0.4.1 (2024-06-13)

### Fix

- fix get name file documentation to fool the release into trueing up

## v0.4.0 (2024-06-13)

### Feat

- update GetFileName function to prefix all file names with an underscore

## v0.3.0 (2024-06-06)

### Fix

- add the implementation context to the GetFileName function

## v0.2.3 (2024-06-05)

### Fix

- merge error fix 2

## v0.2.2 (2024-06-04)

### Fix

- fix merge error

## v0.2.1 (2024-06-04)

### Fix

- add goimports formatting of file to remove unused import files on write

## v0.2.0 (2024-05-31)

### Fix

- remove comment out code from helpers
- update helper functions to be available outside of the package
- update the version file to coincide with go's standard

## 0.1.0 (2024-05-30)

### Feat

- remove CHANGELOG file to get build pipeline working
- rename StructField to Field and update README documentation
- **devops**: added ai workflow for releases
