legendgraphic
=============

`legendgraphic` creates map key/legend for map applications.

You can either use the HTML directly in your web map application, or you can
convert it to a PNG and use it as a legend graphic image in your WMS.


![Example legend](https://github.com/omniscale/legendgraphic/raw/master/example/legend.png)

Installation
------------

Binary releases for various platforms are available at: https://github.com/omniscale/legendgraphic/releases


Configuration
-------------

`legendgraphic` requires a JSON configuration with all layers. Layers are grouped. Each group and each layer can have a title.

Simple legend with one layer:

```json
	{
	    "title": "Legend",
	    "groups": [
		{
		    "title": "Roads and Ways",
		    "layers": [
			{
			    "title": "Motorway",
			    "line-width": 2,
			    "line-color": "#f00",
			    "outline-width": 4,
			    "outline-color": "#ff4444"
			}
		    ]
		}
	    ]
	}
```

Each layer can be symbolized as a line or as a polygon. The following properties define the style of each geometry type.

### Lines

* `line-width`: Width of the line in pixel
* `line-color`: SVG compatible color (e.g. `#ff3322`, `rgb(230, 100, 80)`, `red`)
* `line-dasharray`: Draw line dashed. List of line and spacing intervals in pixels (e.g. `"4, 2, 1, 2"` for a dash-dot-line)
* `outline-width`: Width of the optional outline in pixel 
* `outline-color`: SVG compatible color for outlines

### Polygons

* `fill-color`: SVG compatible color for polygon features

Color variables
---------------

`legendgraphic` supports reading variables from CartoCSS (.mss) files. Instead of defining a color as hex or *rgb* color, you can reference a variable from your CartoCSS files with the @-syntax.

`legendgraphic` accepts one or more mss files as additional parameters:

```
legendgraphic -config config.json -out /tmp/legend.html style.mss
```

Example
-------

```json
{
    "title": "Legend",
    "groups": [
        {
            "title": "Roads and Ways",
            "layers": [
                {
                    "title": "Motorway",
                    "line-width": 2,
                    "line-color": "@motorway_fill",
                    "outline-width": 4,
                    "outline-color": "@motorway_case"
                },
                {
                    "title": "Primary roads",
                    "line-width": 2,
                    "line-color": "@primary_fill",
                    "outline-width": 3,
                    "outline-color": "@primary_case"
                },
                {
                    "title": "Footways",
                    "line-width": 0.5,
                    "line-color": "@footway_fill",
                    "line-dasharray": "@footway_dash"
                }
            ]
        },
        {
            "title": "Land-use",
            "layers": [
                {
                    "title": "Water",
                    "fill-color": "@water"
                },
                {
                    "title": "Park",
                    "fill-color": "@grass"
                },
                {
                    "title": "Forest",
                    "fill-color": "@forest"
                }
            ]
        }
    ]
}
```




Converting to PNG
-----------------

`legendgraphic` only provides HTML with embeded SVG icons. You can create a PNG
by using the screenshot function of your operating system, or you can use
Chrome to automate this.

On MacOS:

    % "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" --headless --disable-gpu --screenshot --window-size=200,600 file:///tmp/legend.html 

