Gyokuro LLC                                           Patrick Coleman 2023/09/1

Title: tinyC2
Author: Patrick Coleman 
Created: September 1st, 2023
Updated: September 28th, 2023

tinyC2 is a simple Command & Control applicaiton. I started this project to
familiarize myself with Golang and Mongo. tinyC2 is true to its name; there
isnt much bells and whistles. IF i manage to find time in the furure I plan
to implement multiple network communication modes such as HTTPS, ICMP and
perhaps some custom protocols.

Below is a diagram of the system components which make tinyC2 possible.
################################################################################
                             ┌───────┐
                             │       │
                             │       │
                             │ mongo │
                             │       │
                             │       │
                             └─┬───▲─┘
                               │   │
                       ┌───────▼───┴───────┐
                       │                   │
                       │    DB Manager     │
                       │                   │
                HTTP   ├──────┬────────────┤
    ┌────────┐   ┌┐    │      │            │
    │        ├───┼┼───►│      │    Node    │
    │   UI   │   ││    │ API  │ Dispatcher │
    │        │◄──┼┼────┤      │            │
    └────────┘   └┘    └──────┘▲───▲───▲──▲┘
                               │   │   │  │
                               │   │   │  │            Golang Channels
┌──────────────────────────────┼───┼───┼──┼───────────────────────────┐
└──────────────────────────────┼───┼───┼──┼───────────────────────────┘
                               │   │   │  │
                               │   │   │  │
                            ┌──┴┬──┴┬──┴┬─┴─┐
                          ┌►│ N │ N │ N │ N │◄┐
                          │ └─▲─┴─▲─┴─▲─┴─▲─┘ │
                          │   │   │   │   │   │
                          │   │   │   │   │   │         CommsPackage Interface
┌─────────────────────────┼───┼───┼───┼───┼───┼───────────────────────┐
└─────────────────────────┼───┼───┼───┼───┼───┼───────────────────────┘
                          │   │   │   │   │   │
                          │   │   │   │   │   │
                        ┌─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┐
                        │ B │ B │ B │ B │ B │ B │
                        └───┴───┴───┴───┴───┴───┘
################################################################################

From the bottom up Beacons (represented by "B" in the diagram) call out to Nodes
(represented by "N" in the diagram). This call occurrs on a periodic adjustable
interval. The communication is handled by a modular interface known as a
"CommsPackage". CommsPackages are designed to be hot-swapped out; they abstract
away the complexities of networking protocols and allow the communications
between beacon and node to be adjusted on the fly. This can be done for a variety
of reasons including to obscure communications or enhance security.

The beacon messages are then transferred to the Node Dispatcher through a shared
channel. The Node Dispatcher is responsible for managing a stable of Nodes and
handling messages between the beacons and the Nodes. The Node Dispatcher will parse
beacon messages, make requests to the DB Manager, and formulate replies for the Node
to send back to the Beacon.

The DB Manager handles all direct database operations. It it utilized by the API
and Node Dispatcher to save the state of the Node Stable as well as record Agent 
Data, Queued Jobs, and log Messages.

The API facilitates the managment of Nodes, Job Queueing, and querying of database
information.

On first start, the Node Dispatcher will query the database for listeners and load
whatever listeners were opeartional prior to shutdown.

################################################################################

Below is the entire flow of registration and completion of the first job.

################################################################################

┌───┐  ┌───┐  ┌───────────────┐  ┌────────────┐  ┌─────┐
│ A │  │ N │  │Node Dispatcher│  │ DB Manager │  │Mongo│
└─┬─┘  └─┬─┘  └───────┬───────┘  └──────┬─────┘  └──┬──┘
  │  1   │            │                 │           │
  ├─────►│     2      │                 │           │
  │      ├───────────►│        3        │           │
  │      │            ├────────────────►│     4     │
  │      │            │                 ├──────────►│
  │      │            │                 │     5     │
  │      │            │        6        │◄──────────┤
  │      │            │◄────────────────│           │
  │      │            │        7        |           │
  │  9   │            │────────────────►│     8     │
  ├─────►│     10     │                 │──────────►│
  │      ├───────────►│        11       │           │
  │      │            ├────────────────►│     12    │
  │      │            │                 ├──────────►│
  │      │            │                 │     13    │
  │      │            │        14       │◄──────────┤
  │      │            │◄────────────────│           │
  │      │            │        15       │           │
  │      │            │────────────────►│     16    │
  │      │            │                 │──────────►│
  │      │            │                 │     17    │
  │      │            │        18       │◄──────────│
  │      │     19     │◄────────────────│           │
  │      │◄───────────│                 │           │
  │  20  │            │                 │           │
  │◄─────|            │                 │           │
  │   21 │            │                 │           │
  │─────►│            │                 │           │
  │      │      22    │                 │           │
  │      │───────────►│                 │           │
  │      │            │        23       │           │
  │      │            │────────────────►│           │
  │      │            │                 │    24     │
  │      │            │                 ├──────────►│
  │      │            │                 │    25     │
  │      │            │                 │◄──────────┤
  │      │            │        26       │           │
  │      │            │◄────────────────│           │
  │      │            │        27       │           │
  │      │            │────────────────►│           │
  │      │            │                 │     28    │
  │      │            │                 │──────────►│
  │      │            │        29       │           │
  │      │            │────────────────►│           │
  │      │            │                 │     30    │
  │      │            │                 │──────────►│
  ▼      ▼            ▼                 ▼           ▼

1. The Agent sends a " hello " message to the Node over a CommsPackage.
2. The Node places this request in its Outbound channel for the Node Dispatcher
    to handle.
3. The Node Dispatcher asks the DB Manager if the Beacon exists.
4. The DB Manager queries the database for the Beacon by name.
5. The database replies with zero results; implying the beacon does not exist.
6. The DB Manager informs the Node Dispatcher that the beacon is new.
7. The Node Dispatcher requests the DB Manager create a new beacon entry.
8. The DB Manager creates a new beacon entry. This marks the end of registration.
9 - 12. This is identical to steps 1 - 4.
13. This time the database returns a result; implying the beacon is registered.
14. The DB Manager replies with one result; implying the beacon does exist.
15. The Node Dispatcher asks the DB Manager if any jobs exist for the beacon.
16. The DB Manager queies the database for open jobs for the beacon.
17. The DB responds with a job.
18. The job is forwarded to the Dispatcher.
19. The job is forwarded to the Node's Inbound Channel.
20. The job is forwarded to the Beacon over a CommsPackage.
21 - 24. This is identical to steps 1 - 4. Except the Beacon's hello message now 
    contains the result of the job from step 20.
25. The database returns a result; implying the beacon is registered.
26. The DB Manager replies with one result; implying the beacon does exist.
27. Node Dispatcher parses the hello message and finds a result from a previous job.
    It requests the DB Manager update the original job request with the result.
28. The DB Manager updates the previous job per step 27.
29. This is identical to step 15. The cycle repeats indefinitely.
################################################################################