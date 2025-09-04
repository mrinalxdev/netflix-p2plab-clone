# My IPLDs

### Making it a Command line interface

I can go with standardized interface like creating a json based interface where external programs describe their DAG structure, the cli tool would execute the external program and parse its JSON output

```json
{
  "nodes": [
    {"id": "node1", "data": {"content": "node1 data"}},
    {"id": "node2", "data": {"content": "node2 data"}}
  ],
  "links": [
    {"source": "node1", "target": "node2", "name": "relation"}
  ]
}
```

more can be done if we go for a package isn't it ?? 

Like devs could import the package, Naah kindaaa shit idea tbh 