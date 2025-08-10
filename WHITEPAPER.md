# Switchyard: An Infrastructure Management Template for Scalable Applications

## Executive Summary

Modern applications that require scalable infrastructure to handle large amounts of user traffic often face increasing levels of complexity in managing distributed systems, feature rollouts, and application reliability. Teams get bogged down with complicated and fragmented tooling when building systems to manage their applications and infrastructure, leading to a high overall system complexity and slower development.

Switchyard addresses these challenges by introducing a unified, plug-and-play infrastructure management toolset that consolidates feature flags, intelligent autoscaling, distributed job scheduling, and incident detection and reporting into a single, cohesive template that deploys on Railway in one click. Being bundled as a Railway template, Switchyard enables teams and developers to deploy infrastructure management capabilities with reduced configuration while maintaining a level of flexibility to customize the management behavior.

## The Problem: Complexity in Existing Infrastructure Management

## Current State Challenges

### Fragmented and Complicated Tooling Options

With the current state of infrastructure management, teams often wrestle many different services and systems that:

-   May not play well together
-   Introduce high levels of complexity
-   Have steep costs

Having to integrate many fragmented services introduces operation complexity, increases costs, and introduces points of failure.

### Autoscaling Limitations

In recent years, Railway has presented itself as an excellent tool for teams and developers looking to deploy applications without having to manage all of the infrastructure themselves. Despite this, Railway lacks some important features that teams look for when building scalable applications, like horizontal scaling capabilities, observability tooling, and robust configuration management. While third party autoscaling integrations exist, they have large fees and may not be flexibile enough for many teams.

## The Solution: Switchyard Platform Architecture

Switchyard provides a unified, pluggable platform that addresses these challenges through several core modules working together:

### Intelligent Autoscaling Engine

Switchyard implements a sophisticated autoscaling algorithm that analyzes and acts on multiple data points to reach scaling decisions:

-   Configurable CPU/Memory thresholds per service
-   Set custom cooldowns for up/downscaling
-   Predictive scaling based on system load
-   Spike and trend analysis for emergency scaling
-   Detection for sustained load
-   Set service replica ranges

### Distributed Job Scheduling and Processing

Switchyard uses a robust job receipt message queuing system with quality-of-service considerations:

-   Requires manual worker acknowledgement for reliability
-   Per-worker concurrency limits to prevent bottlenecking
-   Built-in job idempotency through Redis-based deduplication

Additionally, coupled with the job scaling algorithm, Switchyard autoscales your worker pool to meet system demand.

Teams also define workers as standard Railway services, and Switchyard's architecture allows for unmatched flexibility, enabling:

-   Language agnostic worker implementation
-   Independent deployment and for job logic
-   Horizontal scaling for select job types
-   Isolated failure domains

### Contextual Feature Flag Management

Switchyard's feature flag system supports a targeting rules system, where application services call Switchyard's feature flag server with custom JSON context--based on the rules you set up for each feature flag, Switchyard will evaluate user context against the known set of rules.

Additionally, through Switchyard's web dashboard, feature flag rules and eligibility can be dynamically configured at runtime to respond to system conditions.

### Incident Reporting and Response

Switchyard's incident detection combines multiple sources to provide reporting:

-   Application logs: logs at different severities detected above configurable thresholds will trigger an alarm
-   Interesting service status changes: set different concerning statuses and Switchyard will trigger alerts based on status changes

The incident detection service also provides flexibility by allowing alerts to be sent to custom webhook handlers, allowing:

-   Fanout reporting: send alerts to multiple other services (Slack, PagerDuty, etc.)
-   Incident response: your handler could consume Railway's API or make other requests to respond to incident alerts
-   Chained workflows: handlers can make multiple requests to propogate alerts to other upstream services

### Runtime Configurator

Switchyard also exposes a configurator service that allows for environment variable management at runtime for Switchyard services, enabling:

-   Service observability changes: dynamically change which services are monitored at runtime by changing the incident reporting & locomotive environment variables to add or remove Railway service IDs
-   Adjust monitoring timelines: change monitoring intervals for autoscaling and incident analysis to increase or decrease system response times according to load
-   Adjust error incident reporting thresholds
-   Change application handler URLs at runtime

## Benefits

### Operational Excellence

-   Proactive incident reporting and response handling reduce incident resolution times
-   Intelligent autocaling improve system reliability and reduce service degredation events
-   Cost optimization: compared to traditional licensure and operational fees, Switchyard is open source and self-hosted on Railway

### Development Velocity

-   Single dashboard for all infrastructure concerns reduces context switching and improves operational efficiency.
-   The plug-and-play nature eliminates complex configuration work typically required for complex tools like Kubernetes
-   Switchyard exposes simple APIs that reduce development friction

### Risk Mitigation

-   Feature flags allow for gradual feature rollouts that enable safer deployments and rollbacks
-   Predictive scaling reduces performance incidents and service outages
-   Comprehensive observability: unified monitring and alerting improve system visibility

## Conclusion

Switchyard presents a new solution to fragmented tooling, expensive costs, and a complex developer experience towards a simple, flexible, and robust infrastructure management platform. By consolidating feature flags, autoscaling, distributed job scheduling, and incident detection into a plug-and-play, cohesive system, organizations and teams can achieve higher operating efficiency, improved system reliability, and increased developer velocity.

## Roadmap

In the future, Switchyard will be expanded to support:

-   Configurator support for custom services
-   More feature flag rule operators
-   More considered autoscaling metrics
-   Queue depth analysis for worker pool scaling
