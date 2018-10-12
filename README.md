# here.local
here.local ubiquitous computing software infrastructure


## Components

### Proximity nodes

### Context server

### Here.local server

## Installation

### Here.loca

### Webstrates

## Setting up

## APIs

### Javascript

### Websocket

### RESTful

#### Client information

##### `http://here.local/api/client`
Returns the client device making the request
<pre>
    {
        hostname: "Bobs-macbook",
        mac: "00:00:00:00:00:00:00:00",
        ip: "10.0.0.73",
        vendor: "XEROX CORPORATION",
        agent: "Mozilla/5.0...",
        type: wireless,
        proximity: {
           name: office,
           signal: -60
        },
        locations: [{name:kitchen, signal:-70}, {name:office, signal:-60},...]
    }
</pre>

##### `http://here.local/api/client/<ip>`
Returns client device based on provided IP or empty.
<pre>
    {
        hostname: "Bobs-macbook",
        mac: "00:00:00:00:00:00:00:00",
        ip: "10.0.0.73",
        vendor: "XEROX CORPORATION",
        agent: "Mozilla/5.0...",
        type: wireless,
        signal: -60
    }
</pre>
  
##### `http://here.local/api/client/<mac>`
Returns client device based on provided mac or empty.
<pre>
    {
        hostname: "Bobs-macbook",
        mac: "00:00:00:00:00:00:00:00",
        ip: "10.0.0.73",
        vendor: "XEROX CORPORATION",
        agent: "Mozilla/5.0...",
        type: wireless,
        signal: -60
    }
</pre>

#### Location information


http://here.local/api/location

http://here.local/api/location/<ip>
  
http://here.local/api/location/<mac>

## Limitations and issues
