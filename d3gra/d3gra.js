function draw_charts(id, options) {
    // Размеры и отступы
    let from = ('from' in options) ? options.from : Date.now()-1000;
    let to = ('to' in options) ? options.to : Date.now();
 
    const margin = {top: 4, right: 60, bottom: 20, left: 60};
    const width = (('width' in options) ? options.width : 800) - margin.left - margin.right;
    const height = (('height' in options) ? options.height : 500) - margin.top - margin.bottom;

    // Создание SVG
    const svg = d3.select("#"+id)
    .append("svg")
    .attr("width", width + margin.left + margin.right)
    .attr("height", height + margin.top + margin.bottom)
    .append("g")
    .attr("transform", `translate(${margin.left},${margin.top})`);

    // Шкала для оси X (временная)
    const xScale = d3.scaleTime()
        .domain([from, to])
        .range([0, width]);
    const xAxis = d3.axisBottom(xScale);
    svg.append("g")
        .attr("transform", `translate(0,${height})`)
        .call(xAxis);

    let series = ('series' in options) ? options.series : null;
    var min = 0;
    var max = 0;
    for (var ns = 0; ns < series.length; ns++) {
        var data = series[ns].data;
        for (var nd = 0; nd < data.length; nd++) {
            var dt = data[nd].value;
            if (ns==0 && nd == 0 || dt < min) min = dt;
            if (ns==0 && nd == 0 || dt > max) max = dt;
        }
    }
    if ('min' in options)  {
        let fmin = options.min;
        if (min > fmin) min = fmin;
    }
    if ('max' in options)  {
        let fmax = options.max;
        if (max < fmax) max = fmax;
    }

    // Шкала для оси Y
    var yScale = d3.scaleLinear()
        .domain([min, max])
        .range([height, 0]);
    var yAxis = d3.axisLeft(yScale);
    svg.append("g").call(yAxis);

    const tooltip = d3.select("#tooltip");        
    const say_tip = function(e, th, nm, dt, ds) {
        jQuery("#ttttt").text(nm);
        // Находим ближайшую точку данных
        let [xCoord] = d3.pointer(e, th);
        let bisectDate = d3.bisector(d => d.date).left;
        let x0 = xScale.invert(xCoord);
        let i = bisectDate(dt, x0, 1);
        let dd = dt[i-1];
        // if (i < data.length) {
        //     const d1 = data[i];
        //     if (x0 - dd.date > d1.date - x0) dd = d1;
        // }
        // Позиционируем tooltip
        tooltip.style("opacity", 1)
            .html(`${nm}: ${dd.value}`)
            .style("left", (e.pageX + 10) + "px")
            .style("top", (e.pageY - 28) + "px");
        
        // Подсвечиваем ближайшую точку
        ds.attr("r", 0); // Скрываем все
        d3.select(ds.nodes()[i-1]).attr("r", 5); // Показываем ближайшую
    };
    
    for (var ns = 0; ns < series.length; ns++) {
        var seria = series[ns];
        const data = seria.data;
        const name = ('name' in seria) ? seria.name : ns;
        const clr = ('clr' in seria) ? seria.clr : '#000';

        // Создание генератора линии
        const line = d3.line()
            .x(d => xScale(d.date))
            .y(d => yScale(d.value));
    
        // Рисуем видимую линию
        svg.append("path")
            .datum(data)
            .attr("class", "line")
            .attr("d", line)
            .attr("stroke", clr)
            .attr("fill", "none");
        
        // Создаем НЕВИДИМУЮ область для взаимодействия
        const linePath = svg.append("path")
            .datum(data)
            .attr("class", "line-hover")
            .attr("d", line)
            .attr("stroke-width", 15) // Широкая прозрачная область
            .attr("stroke", "transparent")
            .attr("fill", "none");
        
        // Добавляем точки-маркеры (необязательно)
        const dots = svg.selectAll(".dot-"+ns)
            .data(data)
            .enter()
            .append("circle")
            .attr("class", "dot")
            .attr("cx", d => xScale(d.date))
            .attr("cy", d => yScale(d.value))
            .attr("r", 0) // Изначально невидимы
            .attr("fill", clr);
        
        // var ttp = function(event) {
        //     // Находим ближайшую точку данных
        //     const [xCoord] = d3.pointer(event, this);
        //     const bisectDate = d3.bisector(d => d.date).left;
        //     const x0 = xScale.invert(xCoord);
        //     const i = bisectDate(data, x0, 1);
        //     var dd = data[i-1];
        //     if (i < data.length) {
        //         const d1 = data[i];
        //         if (x0 - dd.date > d1.date - x0) dd = d1;
        //     }
        //     // Позиционируем tooltip
        //     tooltip.style("opacity", 1)
        //         .html(`${name}: ${dd.value}`)
        //         .style("left", (event.pageX + 10) + "px")
        //         .style("top", (event.pageY - 28) + "px");
            
        //     // Подсвечиваем ближайшую точку
        //     dots.attr("r", 0); // Скрываем все
        //     d3.select(dots.nodes()[i-1]).attr("r", 5); // Показываем ближайшую
        // };
        linePath
            // .on("mouseover", function() {
            //     dots.attr("r", 3); // Показываем точки
            // })
            .on("mouseover", function(event){say_tip(event, this, name, data, dots);})
            .on("mousemove", function(event){say_tip(event, this, name, data, dots);})
            .on("mouseout", function() {
                tooltip.style("opacity", 0);
                dots.attr("r", 0); // Скрываем точки
            });
    }
}

