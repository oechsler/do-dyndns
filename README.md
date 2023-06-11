# do-dyndns

A flexible and powerful DynDNS solution powered by DigitalOcean's serverless functions 🚀

## Prerequisites

To embark on this exciting DynDNS adventure, you'll need a couple of things ready:

- A domain in DigitalOcean's free DNS offering. Get ready to bring your domain to life!
- A serverless function namespace in DigitalOcean. It's time to unleash the magic of serverless computing!
- The DigitalOcean CLI installed on your machine.

## Deployment

Now, let's get this show on the road! Here's how you can deploy the do-dyndns project:

Clone the repository to your local machine and navigate to the project directory:

```shell
git clone https://github.com/oechsler/do-dyndns.git
cd do-dyndns
```

Ensure that you have selected the correct serverless function namespace on which you want to deploy the DynDNS functions. Use the following command to connect:

```shell
doctl serverless connect
``` 

It's time to unleash your DynDNS solution to the world!  
Run the deployment command to deploy your code to DigitalOcean's serverless platform:

```shell
doctl serverless deploy .
``` 

Sit back, relax, and watch as your DynDNS solution springs to life. Get ready to experience the magic of dynamic DNS with DigitalOcean's serverless functions!

## Triggering the Functions

All that's missing now is to trigger the functions, for example, from your router's DynDNS service. You can use the `/update4` and `/update6` routes to update either IPv4 or IPv6.
Both routes can be triggered either using query parameters or with the payload in the body of the request.
  
### Parameters
  
| Parameter       | Description                                                                                                                                                                                                                   |
| --------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `record`        | Specifies the subdomain to be updated with the new IP address.                                                                                                                                                       |
| `token`         | Represents the authentication token or password required to authorize the update. Ensure you use a secure and unique token for your setup.                                                                                    |
| `ipv4`          | Used in the IPv4 update route. Specifies the new IPv4 address to associate with the domain or subdomain.                                                                                                                      |
| `ipv6_prefix`   | Used in the IPv6 update route. Represents the prefix of your IPv6 LAN.                                                                                                                                                         |
| `ipv6`          | Used in the IPv6 update route. Specifies the new IPv6 address to associate with the domain or subdomain.                                                                                                                      |
| `interface_id`  | Used in the IPv6 update route. Represents the interface ID for the IPv6 address. Optional parameter, can be provided if needed.                                                                                                |


Now that your do-dyndns project is up and running, you can easily manage and update your DNS records automatically. 
Enjoy the flexibility and convenience of having your domain always point to the right IP address!

Happy DynDNS-ing! 🎉
