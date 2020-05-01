function buildMap(view, points, groups) {
    const groupVectorSource = new ol.source.Vector({
        features: buildGroupsFeatures(groups).concat(buildPointsFeatures(points))
    });

    const clusterSource = new ol.source.Cluster({
        distance: clusterDistance,
        source: groupVectorSource
    });

    const picketStyle = buildSingleStyle();
    const groupStyleCache = {};
    const clusterStyleCache = {};
    const clusters = new ol.layer.Vector({
        source: clusterSource,
        style: function (feature) {
            const groupsAndPoint = feature.get('features');
            if (groupsAndPoint.length === 1) {
                if (groupsAndPoint[0].get("point")) {
                    return picketStyle;
                } else {
                    const group = groupsAndPoint[0].get("group");
                    if (!groupStyleCache[group.count]) {
                        groupStyleCache[group.count] = buildGroupStyle(group.count)
                    }
                    return groupStyleCache[group.count];
                }
            }
            let count = 0;
            groupsAndPoint.forEach((f) => {
                count += f.get('count');
            })
            if (!clusterStyleCache[count]) {
                clusterStyleCache[count] = buildClusterStyle(count)
            }
            return clusterStyleCache[count];
        }
    });

    const map = new ol.Map({
        controls: [
            new ol.control.Zoom()
        ],
        target: 'map',
        layers: [
            new ol.layer.Tile({
                source: new ol.source.OSM(),
            }),
            clusters
        ],
        view: view
    });
    map.on('singleclick', function (evt) {
        clusters.getFeatures(evt.pixel).then(function (clusterFeatures) {
            if (clusterFeatures.length < 1) {
                return;
            }
            view.animate({
                center: clusterFeatures[0].getGeometry().getCoordinates(),
                zoom: view.getZoom() + 1,
                duration: 200
            });
            const features = clusterFeatures[0].get("features");
            if (features.length === 1) {
                if (features[0].get("group")) {
                    showJoinModal(features[0].get('group'));
                }
                if (features[0].get("point")) {
                    showPointModal(features[0].get('point'));
                }
            }
        });
    });

    return map;
}

function buildGroupsFeatures(groups) {
    return groups.map(function (m) {
        const f = new ol.Feature(new ol.geom.Point(ol.proj.fromLonLat([m.longitude, m.latitude])));
        f.set("count", m.count);
        f.set("group", m);
        return f;
    });
}

function buildPointsFeatures(points) {
    return points.map(function (m) {
        const f = new ol.Feature(new ol.geom.Point(ol.proj.fromLonLat([m.longitude, m.latitude])));
        f.set("count", 1);
        f.set("point", m);
        return f;
    });
}

function buildClusterStyle(count) {
    return new ol.style.Style({
        image: new ol.style.Circle({
            radius: radius(count),
            stroke: new ol.style.Stroke({color: '#f14668'}),
            fill: new ol.style.Fill({color: '#feecf0'})
        }),
        text: new ol.style.Text({
            font: "12px sans-serif",
            text: count.toString(),
            fill: new ol.style.Fill({color: '#f14668'})
        }),
        zIndex: count
    });
}

function buildGroupStyle(count) {
    return new ol.style.Style({
        image: new ol.style.Circle({
            radius: radius(count),
            fill: new ol.style.Fill({color: '#f14668'})
        }),
        text: new ol.style.Text({
            font: "12px sans-serif",
            text: count.toString(),
            fill: new ol.style.Fill({color: '#fff'})
        }),
        zIndex: count
    });
}

function buildSingleStyle() {
    return new ol.style.Style({image: new ol.style.Icon({
            opacity: 1,
            src: 'data:image/svg+xml;utf8,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="red" width="24px" height="18px"><path d="M0 0h24v24H0V0z" fill="none"/><path d="M20.75 6.99c-.14-.55-.69-.87-1.24-.75-2.38.53-5.03.76-7.51.76s-5.13-.23-7.51-.76c-.55-.12-1.1.2-1.24.75-.14.56.2 1.13.75 1.26 1.61.36 3.35.61 5 .75v12c0 .55.45 1 1 1s1-.45 1-1v-5h2v5c0 .55.45 1 1 1s1-.45 1-1V9c1.65-.14 3.39-.39 4.99-.75.56-.13.9-.7.76-1.26zM12 6c1.1 0 2-.9 2-2s-.9-2-2-2-2 .9-2 2 .9 2 2 2z"/></svg>',
            scale: 1
        })});
}

function byId(id) {
    return document.getElementById(id);
}

function byClass(cl) {
    return document.getElementsByClassName(cl);
}

function loadData() {
    return fetch('/list', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
        }
    })
        .then((response) => response.json())
        .then((json) => json.data)
        .catch((error) => {
            console.error(error);
        });
}

function radius(count) {
    return 10 + Math.log(count) * 3;
}

function showHelpModal() {
    byId('helpModal').classList.add("is-active");
}

function showJoinModal(data) {
    byId('joinTitle').innerHTML = data.title;
    byId('joinCount').innerHTML = data.count;
    byId('joinText').innerHTML = data.description;
    byId('joinLink').href = "https://t.me/"+data.username;
    byId('joinModal').classList.add("is-active");
}

function showPointModal(data) {
    byId('pointName').innerHTML = data.title;
    byId('pointUsername').innerHTML = "@"+data.username;
    byId('pointText').innerHTML = data.description;
    byId('pointModal').classList.add("is-active");
}

function hideModals() {
    Array.from(byClass("modal")).forEach(element => {
        element.classList.remove("is-active");
    });
}