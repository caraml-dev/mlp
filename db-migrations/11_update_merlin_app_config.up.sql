UPDATE applications
SET config = '{
    "sections": [
        {
            "name": "Models",
            "href": "/models"
        },
        {
            "name": "Standard Transformer Simulator",
            "href": "/simulator"
        },
    ]
}' WHERE name = 'Merlin';