// Code generated by templ@v0.2.364 DO NOT EDIT.

package main

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func render(bm map[string]int, jm map[string]int) templ.ComponentScript {
	return templ.ComponentScript{
		Name: `__templ_render_94b1`,
		Function: `function __templ_render_94b1(bm, jm){google.charts.load("current", { packages: ["corechart", "bar"] });
    google.charts.setOnLoadCallback(drawChart);

    async function drawChart() {
        const data = {
            binary_message_latency_microseconds_bucket: [
                ['Latency (microseconds)', 'Number of Requests'],
                ...Object.entries(bm)
            ],
            json_message_latency_milliseconds_bucket: [
                ['Latency (milliseconds)', 'Number of Requests'],
                ...Object.entries(jm)
            ]
        };

        Object.entries(data).forEach(([key, value]) => {
            const data = google.visualization.arrayToDataTable(value);

            const options = {
                title: key.replace("_", " ").replace("bucket", ""),
                legend: { position: "none" },
            };

            const body = document.querySelector("body");

            const div = document.createElement("div");

            div.setAttribute("style", "width: 900px; height: 500px; margin-bottom: 20px")

            body.appendChild(div);

            const chart = new google.charts.Bar(div);

            chart.draw(data, options);
        });
    }}`,
		Call: templ.SafeScript(`__templ_render_94b1`, bm, jm),
	}
}

func chart(bm map[string]int, jm map[string]int) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_1 := templ.GetChildren(ctx)
		if var_1 == nil {
			var_1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<html><head><script type=\"text/javascript\" src=\"https://www.gstatic.com/charts/loader.js\">")
		if err != nil {
			return err
		}
		var_2 := ``
		_, err = templBuffer.WriteString(var_2)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</script></head>")
		if err != nil {
			return err
		}
		err = templ.RenderScriptItems(ctx, templBuffer, render(bm, jm))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<body onload=\"")
		if err != nil {
			return err
		}
		var var_3 templ.ComponentScript = render(bm, jm)
		_, err = templBuffer.WriteString(var_3.Call)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"></body></html>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}
