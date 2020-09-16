# Concepts

## Machine Learning Platform

Machine Learning Platform (MLP) aims to solve the following problems:

- The data science development experience can be painful
- No standardisation of the machine learning life cycle
- Duplicated effort and lack of reusability
- Difficult to get data science systems into production
- Hard to maintain data science systems once in production
- Difficult to measure impact

MLP's vision is to empower data scientists, analysts, and other ML developers to create ML solutions that drive direct business impact. These solutions can range from simple analyses all the way to production ML systems that serve millions of customers. The ML Platform aims to provide these users with a unified set of products they can use to rapidly develop and confidently deploy their ML solutions.

## Machine Learning Life Cycle

The typical ML life cycle can be viewed through the following nine stages:

![Machine learning life cycle](./diagrams/machine_learning_life_cycle.drawio.svg)

Starting by (1) **sourcing data**, a data scientist will (2) **explore and analyze** it. The raw data is (3) **transformed** into useful features, typically involving (4) **scheduling and automation** to do this on a regular basis. The resultant features are (5) **stored and managed**, available for the various models and other data scientists to use. As part of the exploration, the data scientist will also (6) **build, train, and evaluate** various models. Promising models are (7) **stored and deployed** into production. The production models are then (8) **served and monitored** for a period of time. Typically, there are multiple competing models in production, and choosing between them or evaluating them is done via (9) **experimentation**. With the learnings of the production models, the data scientist iterates on new features and models.

## MLP Product

MLP Products are systems and services that are specifically built to solve one or multiple stages of the machine learning life cycle's problems.

MLP Products share the following design principles:

- **Easy to compose ML solutions out of parts of the platform** - New [ML projects](#ml-project) should be able to compose solutions out of existing MLP Products on the ML Platform, instead of building from scratch. With the infrastructure complexity abstracted away, the entry barrier of using ML to drive business impact is lowered and would allow a lightweight data science team or even non-data scientists to leverage ML power.
- **Best practices are enforced and unified on each stage in the machine learning life cycle** - Data scientists should have a clear understanding of all the stages of the ML life cycle, the products that exist at each stage, and how to apply them to their use cases in a self-service manner with minimal support from the engineers.
- **Integration into the existing tech stack** - The MLP Product is compatible with the existing tech stack and either abstracts away any integration points or makes these integrations easy. Data scientists should not have to be concerned with how their solutions will be consumed. Furthermore, the platform leverages many of the existing products and tools provided by the open-source communities.
- **Bottom-up innovation** - The platform is built in a modular fashion, in layers from the ground up. Given the diversity of use cases and applications that need to be supported, it is necessary to support not only the “happy path”, but to also provide flexibility when edge cases arise.

Currently, we have published the following MLP Products:

* [**Feast**](https://github.com/gojek/feast) - For managing and serving machine learning features.
* [**Merlin**](https://github.com/gojek/merlin) - For deploying, serving, and monitoring machine learning models.
* [**Turing**](https://github.com/gojek/turing) - For designing, deploying, and evaluating machine learning experiments.

## ML Project

Machine Learning Project organizes all your ML resources and solutions and links them together. An ML Project can consists of, including, but not limited to, raw data, a set of features, models, and experimentations. For example, an ML Project could be a driver allocation system consists of pipelines for data transformation, feature engineering, and models training; and models serving and experimentation in production.
