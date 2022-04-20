# PiHatDraw
The PiHatDraw is a learning project, to create a drwaing application to run with a Raspberry Pi with a
[Sense HAT](https://www.raspberrypi.org/products/sense-hat/), in the go programming language (goland). 

The Sense HAT contains an 8X8 LED display and a joystick. In this project, we're using the joystick to draw. 
The Sense HAT LED display is a floating window that shows a subset of the full picture. 

![Sense HAT](readmeImages/pi-hat.jpeg)

The full picture can be viewed using a web browser. The web display also contains some controllers, like changing the pen color and so on.

![Web Application display](readmeImages/weapp.png)

The application is written in golang. It's an event driven application. It uses webSockets to keep the web display in sync, and a web application to handle requests from the web client.

I wrote a blog post series with a tutorial to guide how to build this application.

For each post, there is a coresponding tag with the relevant code:
* [Introduction post](https://nunnatsa.github.io/piHatDraw/)
* [Chapter 1: Start Drawing](https://nunnatsa.github.io/piHatDraw/ch1.html) - [v0.0.1](https://github.com/nunnatsa/piHatDraw/releases/tag/v0.0.1)
* [Chapter 2: Add Web-based display](https://nunnatsa.github.io/piHatDraw/ch2.html) - [v0.0.2](https://github.com/nunnatsa/piHatDraw/releases/tag/v0.0.2)
* [Chapter 3: Expand Canvas](https://nunnatsa.github.io/piHatDraw/ch3.html) - [v0.0.3](https://github.com/nunnatsa/piHatDraw/releases/tag/v0.0.3)
* [Chapter 4: Adding Colors](https://nunnatsa.github.io/piHatDraw/ch4.html) - [v0.0.4](https://github.com/nunnatsa/piHatDraw/releases/tag/v0.0.4)
* [Chapter 5: Download the Picture](https://nunnatsa.github.io/piHatDraw/ch5.html) - [v0.0.5](https://github.com/nunnatsa/piHatDraw/releases/tag/v0.0.5)
* [Chapter 6: Undo](https://nunnatsa.github.io/piHatDraw/ch6.html) - [v0.0.6](https://github.com/nunnatsa/piHatDraw/releases/tag/v0.0.6)
* [Chapter 7: The Bucket Tool](https://nunnatsa.github.io/piHatDraw/ch7.html) - [v0.0.7](https://github.com/nunnatsa/piHatDraw/releases/tag/v0.0.7)
* Working on the next chapter: Replacing the UI to vue and vuetify, and compiling for another architecture.

## Demo
[<img src="https://i3.ytimg.com/vi/2IngYHPHjtc/maxresdefault.jpg" width="50%">](https://youtu.be/2IngYHPHjtc "click for video with the demo")