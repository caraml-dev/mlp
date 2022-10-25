UPDATE applications
SET config = '{
    "sections": [
        {
            "name": "Models",
            "href": "/models"
        },
        {
            "name": "Standard Transformer Simulator",
            "href": "/transformer-simulator"
        }
    ]
}' WHERE name = 'Merlin';