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

```
    {
        hostname: "Bobs-macbook",
        mac: "00:00:00:00:00:00:00:00",
        ip: "10.0.0.73",
        vendor: "XEROX CORPORATION",
        agent: "Mozilla/5.0...",
        signal: -60,
        proximity: {
            
        }
    }
```

##### `http://here.local/api/client/<ip>`
<pre>
    {
        hostname: "Bobs-macbook",
        mac: "00:00:00:00:00:00:00:00",
        <b>ip: "10.0.0.73"</b>,
        vendor: "XEROX CORPORATION",
        agent: "Mozilla/5.0...",
        signal: -60,
        proximity: {
            
        }
    }
</pre>
  
##### `http://here.local/api/client/<mac>`


#### Location information


http://here.local/api/location

http://here.local/api/location/<ip>
  
http://here.local/api/location/<mac>

## Limitations and issues
