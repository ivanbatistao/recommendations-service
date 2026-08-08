# Problem Statement
Everything start with the following question:
- What problem are we solving?

## Context
An e-commerce platform records millions of user interactions every day:
- Product views
- Searches
- Add-to-cart events
- Purchases
- Each interaction contains valuable information that can be used to generate personalized recommendations.

Currently, when a user requests recommendations, the system must query multiple data sources and calculate the recommendations at request time, resulting in high latency.

## Problem
The response time for retrieving recommendations is too high, negatively affecting the user experience and reducing the likelihood of a purchase.

## Objective
Build a system that continuously processes events and maintains precomputed recommendations so that requests can be served with low latency