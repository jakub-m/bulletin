# Usage tips

Paste some potential feeds to `input.txt` and run

```bash
cat input.txt | xargs -n1 ./bin/bulletin test  >> feeds.conf
cat feeds.conf > t ; cat t | sort | uniq > feeds.conf; rm -f t
```
