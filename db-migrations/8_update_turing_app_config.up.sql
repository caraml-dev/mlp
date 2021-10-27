UPDATE applications
SET config = '{
    "sections": [
        {
            "name": "Routers",
            "href": "/routers"
        },
        {
            "name": "Ensemblers",
            "href": "/ensemblers"
        },
        {
            "name": "Ensembling Jobs",
            "href": "/jobs"
        },
        {
            "name": "Experiments",
            "href": "/experiments"
        }
    ]
}' WHERE name = 'Turing';