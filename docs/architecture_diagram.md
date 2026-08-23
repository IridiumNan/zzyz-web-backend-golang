# System Architecture Diagram

```mermaid
graph TB
    subgraph zzyz_web
    raw_dir_zip --> |unzip|markdown_pictures_dir
    markdown_pictures_dir --> |parse_toml_config|database_metadata
    database --> |present|admin --> |pass|build_html
    end
    member --> |upload_zip|goalkeeper
    goalkeeper --> |upzip_and_sort|package
    package --> zzyz_web
    user --> |request|zzyz_web
    zzyz_web --> |post_list_html|user
```

```mermaid
graph TB
    subgraph PostList
    post_overview_1
    post_overview_2
    post_overview_3
    end

    subgraph post_overview_1
    title
    short_information
    html_link
    end
```
