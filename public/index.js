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

function getRandomInt(max) {
    return Math.floor(Math.random() * Math.floor(max));
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
    byId('joinDesc').innerHTML = data.description.replace(/(?:\r\n|\r|\n)/g, '<br>');
    byId('joinLink').href = data.link;
    byId('joinModal').classList.add("is-active");
}

function showPointModal(data) {
    console.log(data);
    byId('pointDesc').innerHTML = data.description.replace(/(?:\r\n|\r|\n)/g, '<br>');
    byId('pointImage').src = data.photo.Path;
    byId('pointModal').classList.add("is-active");
}

function hideModals() {
    Array.from(byClass("modal")).forEach(element => {
        element.classList.remove("is-active");
    });
}