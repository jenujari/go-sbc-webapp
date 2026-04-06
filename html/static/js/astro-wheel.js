(() => {
    const root = document.currentScript?.previousElementSibling;
    const svgNode = root?.querySelector("#astro-wheel");
    const tooltip = root?.querySelector("#astro-wheel-tooltip");
    const wheelWrap = svgNode?.parentElement;
    const planetsDataInput = root?.querySelector("#planets-data");
    const planets = JSON.parse(planetsDataInput.value) || [];

    if (!root || !svgNode || !wheelWrap || typeof d3 === "undefined") return;

    const zodiac = [
      "Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo",
      "Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces"
    ];
    const nakshatraSegments = [
      { label: "Ashwini", start: 0, end: 13.333333 },
      { label: "Bharani", start: 13.333333, end: 26.666666 },
      { label: "Krittika", start: 26.666666, end: 40 },
      { label: "Rohini", start: 40, end: 53.333333 },
      { label: "Mrigashirsha", start: 53.333333, end: 66.666666 },
      { label: "Ardra", start: 66.666666, end: 80 },
      { label: "Punarvasu", start: 80, end: 93.333333 },
      { label: "Pushya", start: 93.333333, end: 106.666666 },
      { label: "Ashlesha", start: 106.666666, end: 120 },
      { label: "Magha", start: 120, end: 133.333333 },
      { label: "Purva Phalguni", start: 133.333333, end: 146.666666 },
      { label: "Uttara Phalguni", start: 146.666666, end: 160 },
      { label: "Hasta", start: 160, end: 173.333333 },
      { label: "Chitra", start: 173.333333, end: 186.666666 },
      { label: "Swati", start: 186.666666, end: 200 },
      { label: "Vishakha", start: 200, end: 213.333333 },
      { label: "Anuradha", start: 213.333333, end: 226.666666 },
      { label: "Jyeshtha", start: 226.666666, end: 240 },
      { label: "Moola", start: 240, end: 253.333333 },
      { label: "Purva Ashadha", start: 253.333333, end: 266.666666 },
      { label: "Uttara Ashadha", start: 266.666666, end: 276.666666 },
      { label: "Abhijit", start: 276.666666, end: 280.888889 },
      { label: "Shravana", start: 280.888889, end: 293.333333 },
      { label: "Dhanishtha", start: 293.333333, end: 306.666666 },
      { label: "Shatabhisha", start: 306.666666, end: 320 },
      { label: "Purva Bhadrapada", start: 320, end: 333.333333 },
      { label: "Uttara Bhadrapada", start: 333.333333, end: 346.666666 },
      { label: "Revati", start: 346.666666, end: 360 }
    ];
    const planetColors = {
      sun: "#f59e0b",
      moon: "#38bdf8",
      mercury: "#10b981",
      venus: "#ec4899",
      mars: "#ef4444",
      jupiter: "#8b5cf6",
      saturn: "#334155",
      uranus: "#06b6d4",
      neptune: "#6366f1",
      pluto: "#7c2d12",
      rahu: "#0f766e",
      ketu: "#a16207"
    };
    const signFills = ["#fde68a", "#bfdbfe", "#fbcfe8"];
    const size = 920;
    const center = size / 2;
    const orbitGuide = 430;
    const signInner = 205;
    const signOuter = 300;
    const nakInner = 310;
    const nakOuter = 372;
    const solarOrbitRadius = 404;
    const earthRadius = 76;
    const clusterThreshold = 6;
    const arcGen = d3.arc();
    const sunPlanet = planets.find((planet) => planet.Name.toLowerCase() === "sun") || null;
    const sunLatitude = sunPlanet ? sunPlanet.Lat : 0;
    const relativeLatitudes = planets.map((planet) => planet.Lat - sunLatitude);
    const maxRelativeMagnitude = relativeLatitudes.length
      ? Math.max(...relativeLatitudes.map((value) => Math.abs(value))) || 1
      : 1;

    const degToRad = (degrees) => (degrees * Math.PI) / 180;
    const toArcAngle = (degrees) => ((90 - degrees) * Math.PI) / 180;
    const normalizeAngle = (degrees) => ((degrees % 360) + 360) % 360;
    const shortestAngleDiff = (a, b) => {
      const diff = Math.abs(normalizeAngle(a) - normalizeAngle(b));
      return Math.min(diff, 360 - diff);
    };
    const pointAt = (radius, degrees) => ({
      x: Math.cos(degToRad(degrees)) * radius,
      y: -Math.sin(degToRad(degrees)) * radius
    });

    const clusteredPlanets = [...planets]
      .sort((a, b) => a.Long - b.Long)
      .map((planet) => ({ ...planet, clusterIndex: 0, clusterSize: 1 }));

    const clusters = [];
    clusteredPlanets.forEach((planet) => {
      const lastCluster = clusters[clusters.length - 1];
      if (!lastCluster || shortestAngleDiff(lastCluster[lastCluster.length - 1].Long, planet.Long) > clusterThreshold) {
        clusters.push([planet]);
        return;
      }
      lastCluster.push(planet);
    });

    if (clusters.length > 1) {
      const firstCluster = clusters[0];
      const lastCluster = clusters[clusters.length - 1];
      if (shortestAngleDiff(firstCluster[0].Long, lastCluster[lastCluster.length - 1].Long) <= clusterThreshold) {
        clusters[0] = lastCluster.concat(firstCluster);
        clusters.pop();
      }
    }

    clusters.forEach((cluster) => {
      cluster.forEach((planet, index) => {
        planet.clusterIndex = index;
        planet.clusterSize = cluster.length;
      });
    });

    const renderWheel = () => {
      const svg = d3.select(svgNode);
      svg.selectAll("*").remove();
      tooltip?.classList.add("hidden");
      if (tooltip) tooltip.style.display = "none";

      svg
        .attr("viewBox", `0 0 ${size} ${size}`)
        .attr("preserveAspectRatio", "xMidYMid meet");

      svg.append("rect")
        .attr("width", size)
        .attr("height", size)
        .attr("fill", "#ffffff");

      const g = svg.append("g")
        .attr("transform", `translate(${center}, ${center})`);

      g.append("circle")
        .attr("r", orbitGuide)
        .attr("fill", "#f8fafc")
        .attr("stroke", "#e2e8f0");

      g.append("circle")
        .attr("r", signInner - 10)
        .attr("fill", "#e0f2fe")
        .attr("stroke", "#cbd5e1");

      g.append("circle")
        .attr("r", solarOrbitRadius)
        .attr("fill", "none")
        .attr("stroke", "#f59e0b")
        .attr("stroke-width", 2.5)
        .attr("stroke-dasharray", "8 7")
        .attr("opacity", 0.85);

      g.selectAll(".zodiac-arc")
        .data(zodiac)
        .enter()
        .append("path")
        .attr("class", "arc zodiac-arc")
        .attr("d", (d, index) => arcGen({
          innerRadius: signInner,
          outerRadius: signOuter,
          startAngle: toArcAngle((index + 1) * 30),
          endAngle: toArcAngle(index * 30)
        }))
        .attr("fill", (d, index) => signFills[index % signFills.length]);

      g.selectAll(".zodiac-label")
        .data(zodiac)
        .enter()
        .append("text")
        .attr("transform", (d, index) => {
          const point = pointAt((signInner + signOuter) / 2, index * 30 + 15);
          return `translate(${point.x}, ${point.y})`;
        })
        .attr("font-size", 17)
        .attr("fill", "#334155")
        .attr("font-weight", 700)
        .text((d) => d);

      g.selectAll(".nakshatra-arc")
        .data(nakshatraSegments)
        .enter()
        .append("path")
        .attr("class", "nakshatra nakshatra-arc")
        .attr("d", (segment) => arcGen({
          innerRadius: nakInner,
          outerRadius: nakOuter,
          startAngle: toArcAngle(segment.end),
          endAngle: toArcAngle(segment.start)
        }))
        .attr("fill", (d, index) => index % 2 ? "#dbeafe" : "#fef3c7");

      g.selectAll(".nakshatra-label")
        .data(nakshatraSegments)
        .enter()
        .append("text")
        .attr("transform", (segment) => {
          const point = pointAt((nakInner + nakOuter) / 2, (segment.start + segment.end) / 2);
          return `translate(${point.x}, ${point.y})`;
        })
        .attr("font-size", (segment) => segment.label.length > 12 ? 10 : 11)
        .attr("font-weight", 600)
        .attr("fill", "#7c2d12")
        .text((segment) => segment.label);

      const degrees = d3.range(0, 360, 30);

      g.selectAll(".deg-line")
        .data(degrees)
        .enter()
        .append("line")
        .attr("class", "degree-line deg-line")
        .attr("x1", (degree) => pointAt(signInner, degree).x)
        .attr("y1", (degree) => pointAt(signInner, degree).y)
        .attr("x2", (degree) => pointAt(orbitGuide, degree).x)
        .attr("y2", (degree) => pointAt(orbitGuide, degree).y);

      g.selectAll(".deg-label")
        .data(degrees)
        .enter()
        .append("text")
        .attr("x", (degree) => pointAt(signInner - 34, degree).x)
        .attr("y", (degree) => pointAt(signInner - 34, degree).y)
        .attr("font-size", 14)
        .attr("fill", "#475569")
        .attr("font-weight", 700)
        .text((degree) => `${degree}°`);

      [signInner, signOuter, nakInner, nakOuter].forEach((radius) => {
        g.append("circle")
          .attr("r", radius)
          .attr("fill", "none")
          .attr("stroke", "#cbd5e1")
          .attr("stroke-width", 1.5);
      });

      const showTooltip = (event, planet, color) => {
        if (!tooltip) return;
        const relativeLatitude = planet.Lat - sunLatitude;
        tooltip.innerHTML = `
          <div class="font-semibold font-noto-serif text-slate-900">${planet.Name}</div>
          <div class="mt-1 text-xs font-medium font-roboto text-slate-600">Longitude: ${planet.Longitude}</div>
          <div class="text-xs font-medium font-roboto text-slate-500">Latitude: ${planet.Latitude}</div>
          <div class="text-xs font-medium font-roboto text-slate-500">Relative to Sun: ${relativeLatitude >= 0 ? "+" : ""}${relativeLatitude.toFixed(2)}°</div>
        `;
        tooltip.style.borderColor = color;
        tooltip.style.display = "block";
        tooltip.classList.remove("hidden");

        const bounds = wheelWrap.getBoundingClientRect();
        tooltip.style.left = `${event.clientX - bounds.left + 14}px`;
        tooltip.style.top = `${event.clientY - bounds.top - 12}px`;
      };

      const hideTooltip = () => {
        if (!tooltip) return;
        tooltip.style.display = "none";
        tooltip.classList.add("hidden");
      };

      const planetGroup = g.selectAll(".planet")
        .data(clusteredPlanets)
        .enter()
        .append("g")
        .attr("class", "planet")
        .attr("cursor", "pointer")
        .each(function (planet) {
          const relativeLatitude = planet.Lat - sunLatitude;
          const normalizedRelativeLatitude = relativeLatitude / maxRelativeMagnitude;
          const clusterCenter = (planet.clusterSize - 1) / 2;
          const tangentialOffset = (planet.clusterIndex - clusterCenter) * 18;
          const clusterRadialOffset = relativeLatitude === 0
            ? (planet.clusterIndex - clusterCenter) * 5
            : Math.sign(relativeLatitude) * Math.abs(planet.clusterIndex - clusterCenter) * 6;
          const isSun = planet.Name.toLowerCase() === "sun";
          const planetRadius = isSun
            ? solarOrbitRadius
            : solarOrbitRadius + normalizedRelativeLatitude * 30 + clusterRadialOffset;
          const angle = normalizeAngle(planet.Long + (tangentialOffset / planetRadius) * (180 / Math.PI));
          const anchor = pointAt(isSun ? solarOrbitRadius : nakOuter + 3, angle);
          const position = pointAt(planetRadius, angle);
          const color = planetColors[planet.Name.toLowerCase()] || "#2563eb";
          const shortName = planet.Name.length > 2 ? planet.Name.slice(0, 2).toUpperCase() : planet.Name.toUpperCase();

          planet.render = { anchor, position, color, shortName };
        });

      planetGroup.append("circle")
        .attr("r", 17)
        .attr("cx", (planet) => planet.render.position.x)
        .attr("cy", (planet) => planet.render.position.y)
        .attr("fill", (planet) => planet.render.color)
        .attr("opacity", 0);

      planetGroup.append("line")
        .attr("x1", (planet) => planet.render.anchor.x)
        .attr("y1", (planet) => planet.render.anchor.y)
        .attr("x2", (planet) => planet.render.position.x)
        .attr("y2", (planet) => planet.render.position.y)
        .attr("stroke", (planet) => planet.render.color)
        .attr("stroke-width", 2)
        .attr("stroke-linecap", "round")
        .attr("opacity", 0.9);

      planetGroup.append("circle")
        .attr("r", 10.5)
        .attr("cx", (planet) => planet.render.position.x)
        .attr("cy", (planet) => planet.render.position.y)
        .attr("fill", (planet) => planet.render.color)
        .attr("stroke", "#ffffff")
        .attr("stroke-width", 2);

      planetGroup.append("text")
        .attr("class", "planet-code")
        .attr("x", (planet) => planet.render.position.x)
        .attr("y", (planet) => planet.render.position.y + 0.5)
        .text((planet) => planet.render.shortName);

      planetGroup
        .on("mouseenter", function (event, planet) {
          const group = d3.select(this);
          group.select("circle").attr("opacity", 0.18);
          group.select("line").attr("stroke-width", 3.5);
          group.selectAll("circle").filter((d, index) => index === 1).attr("r", 12.5);
          showTooltip(event, planet, planet.render.color);
        })
        .on("mousemove", function (event, planet) {
          showTooltip(event, planet, planet.render.color);
        })
        .on("mouseleave", function () {
          const group = d3.select(this);
          group.select("circle").attr("opacity", 0);
          group.select("line").attr("stroke-width", 2);
          group.selectAll("circle").filter((d, index) => index === 1).attr("r", 10.5);
          hideTooltip();
        });

      g.append("circle")
        .attr("r", earthRadius)
        .attr("fill", "#e0f2fe")
        .attr("stroke", "#94a3b8")
        .attr("stroke-width", 2);

      g.append("text")
        .attr("y", -8)
        .attr("fill", "#0f172a")
        .attr("font-size", 30)
        .attr("font-weight", 700)
        .text("Earth");

      g.append("text")
        .attr("y", 26)
        .attr("fill", "#64748b")
        .attr("font-size", 14)
        .attr("font-weight", 500)
        .text("Centered Reference");
    };

    requestAnimationFrame(renderWheel);
  })();