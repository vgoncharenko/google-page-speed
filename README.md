For using Google Page Speed test you need to create little script for:
1) testing some pages for your site - it's like some scenarios
2) sending requests for each scenario
3) analyzing response (in error case or inconsistency response format)  for each scenario
4) store results for each scenarios
5) And store the main data 'SCORE' from responses of each scenarios to one report file

You can do all of this by next line:
`/usr/bin/googlepagespeed
                -url=${request_protocol}://${host_name}
                -apiKey=${test.googlepagespeed.api_key}
                -scenarios=/path/config.json'`

/path/config.json have next format:
```
[
    {"scenario_name":"Home page", "sub_url":""},
    {"scenario_name":"Category page", "sub_url":"category-1.html"},
    {"scenario_name":"Product page", "sub_url":"simple-product-4.html"}
]
```

All results will be stored in /root/googlepagespeed/
docker-compose.yml:
```
services:
  googlepagespeed:
    image: lord723518/google-page-speed
    entrypoint: tail -f /dev/null
    working_dir: /root/googlepagespeed
    volumes:
      - ./googlepagespeed:/root/googlepagespeed:rw
```