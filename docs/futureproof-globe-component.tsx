import { useState, useEffect, useRef, useCallback } from "react";
import * as THREE from "three";

// ── DATA ─────────────────────────────────────────────────────────────────────
const CITIES = [
  { id:1,  name:"San Francisco", country:"USA",          region:"Americas", lat:37.77,  lng:-122.41, score:95, growth:12, remote:85, salaryIdx:98, jobs:["AI Engineer","Quantum Developer","Prompt Engineer"], color:0x00ff88 },
  { id:2,  name:"New York",      country:"USA",          region:"Americas", lat:40.71,  lng:-74.00,  score:91, growth:9,  remote:78, salaryIdx:94, jobs:["FinTech Architect","AI Ethics Officer","Cybersecurity Lead"], color:0x00ff88 },
  { id:3,  name:"Austin",        country:"USA",          region:"Americas", lat:30.26,  lng:-97.74,  score:89, growth:14, remote:83, salaryIdx:88, jobs:["Space Tech Engineer","AI Developer","Clean Energy Lead"], color:0x00aaff },
  { id:4,  name:"Seattle",       country:"USA",          region:"Americas", lat:47.60,  lng:-122.33, score:90, growth:11, remote:82, salaryIdx:95, jobs:["Cloud Architect","AI Researcher","Quantum Developer"], color:0x00ff88 },
  { id:5,  name:"Toronto",       country:"Canada",       region:"Americas", lat:43.65,  lng:-79.38,  score:86, growth:13, remote:82, salaryIdx:82, jobs:["AI Developer","Climate Tech Engineer","Biotech Researcher"], color:0x00aaff },
  { id:6,  name:"Vancouver",     country:"Canada",       region:"Americas", lat:49.28,  lng:-123.12, score:85, growth:12, remote:85, salaryIdx:80, jobs:["AI Engineer","Climate Tech Lead","Digital Media Creator"], color:0x00aaff },
  { id:7,  name:"São Paulo",     country:"Brazil",       region:"Americas", lat:-23.55, lng:-46.63,  score:72, growth:16, remote:68, salaryIdx:48, jobs:["AgriTech Developer","FinTech Lead","Climate Engineer"], color:0xffaa00 },
  { id:8,  name:"London",        country:"UK",           region:"Europe",   lat:51.50,  lng:-0.12,   score:88, growth:10, remote:80, salaryIdx:90, jobs:["AI Ethics Officer","Green Energy Consultant","BioTech Researcher"], color:0x00aaff },
  { id:9,  name:"Berlin",        country:"Germany",      region:"Europe",   lat:52.52,  lng:13.40,   score:85, growth:11, remote:78, salaryIdx:78, jobs:["Sustainability Engineer","Robotics Developer","Cybersecurity Analyst"], color:0x00aaff },
  { id:10, name:"Stockholm",     country:"Sweden",       region:"Europe",   lat:59.33,  lng:18.07,   score:89, growth:14, remote:85, salaryIdx:85, jobs:["Green Tech Engineer","AI Ethics Lead","Digital Wellness Coach"], color:0x00aaff },
  { id:11, name:"Amsterdam",     country:"Netherlands",  region:"Europe",   lat:52.37,  lng:4.89,    score:87, growth:12, remote:83, salaryIdx:84, jobs:["Sustainability Architect","Data Engineer","Cybersecurity Expert"], color:0x00aaff },
  { id:12, name:"Zurich",        country:"Switzerland",  region:"Europe",   lat:47.37,  lng:8.54,    score:90, growth:10, remote:82, salaryIdx:100,jobs:["FinTech Architect","Biotech Scientist","Quantum Developer"], color:0x00ff88 },
  { id:13, name:"Helsinki",      country:"Finland",      region:"Europe",   lat:60.17,  lng:24.94,   score:88, growth:13, remote:86, salaryIdx:83, jobs:["EdTech Developer","AI Researcher","Green Energy Specialist"], color:0x00aaff },
  { id:14, name:"Paris",         country:"France",       region:"Europe",   lat:48.85,  lng:2.35,    score:84, growth:9,  remote:76, salaryIdx:82, jobs:["AI Researcher","Digital Arts Lead","BioTech Developer"], color:0x00aaff },
  { id:15, name:"Warsaw",        country:"Poland",       region:"Europe",   lat:52.23,  lng:21.01,   score:79, growth:17, remote:75, salaryIdx:62, jobs:["Software Engineer","Cybersecurity Analyst","AI Developer"], color:0xffaa00 },
  { id:16, name:"Lisbon",        country:"Portugal",     region:"Europe",   lat:38.72,  lng:-9.14,   score:78, growth:15, remote:82, salaryIdx:60, jobs:["Remote Work Hub Lead","Digital Nomad Specialist","Green Tech Engineer"], color:0xffaa00 },
  { id:17, name:"Tel Aviv",      country:"Israel",       region:"Middle East",lat:32.08, lng:34.78,  score:91, growth:17, remote:80, salaryIdx:88, jobs:["Cybersecurity Expert","AI Researcher","BioTech Developer"], color:0x00ff88 },
  { id:18, name:"Dubai",         country:"UAE",          region:"Middle East",lat:25.20, lng:55.27,  score:82, growth:18, remote:65, salaryIdx:85, jobs:["Smart City Planner","Space Tourism Specialist","Digital Finance Expert"], color:0x00aaff },
  { id:19, name:"Singapore",     country:"Singapore",    region:"Asia",     lat:1.35,   lng:103.82,  score:92, growth:15, remote:75, salaryIdx:91, jobs:["FinTech Developer","Smart City Architect","Climate Tech Lead"], color:0x00ff88 },
  { id:20, name:"Tokyo",         country:"Japan",        region:"Asia",     lat:35.68,  lng:139.69,  score:87, growth:9,  remote:60, salaryIdx:80, jobs:["Robotics Engineer","Neural Interface Dev","Space Tech Specialist"], color:0x00aaff },
  { id:21, name:"Seoul",         country:"South Korea",  region:"Asia",     lat:37.57,  lng:126.97,  score:84, growth:11, remote:65, salaryIdx:76, jobs:["Metaverse Developer","Robotics Engineer","Semiconductor Architect"], color:0x00aaff },
  { id:22, name:"Beijing",       country:"China",        region:"Asia",     lat:39.90,  lng:116.40,  score:86, growth:15, remote:50, salaryIdx:72, jobs:["AI Engineer","Robotics Developer","Space Tech Expert"], color:0x00aaff },
  { id:23, name:"Bangalore",     country:"India",        region:"Asia",     lat:12.97,  lng:77.59,   score:83, growth:20, remote:70, salaryIdx:55, jobs:["AI Engineer","Quantum Computing Dev","Space Tech Specialist"], color:0x00aaff },
  { id:24, name:"Shenzhen",      country:"China",        region:"Asia",     lat:22.54,  lng:114.06,  score:85, growth:16, remote:52, salaryIdx:74, jobs:["Hardware Engineer","AI Developer","Smart Manufacturing Lead"], color:0x00aaff },
  { id:25, name:"Sydney",        country:"Australia",    region:"Oceania",  lat:-33.86, lng:151.20,  score:80, growth:8,  remote:79, salaryIdx:83, jobs:["Climate Scientist","Mining Tech Engineer","Digital Health Expert"], color:0xffaa00 },
  { id:26, name:"Melbourne",     country:"Australia",    region:"Oceania",  lat:-37.81, lng:144.96,  score:79, growth:9,  remote:80, salaryIdx:81, jobs:["BioTech Researcher","AgriTech Engineer","AI Developer"], color:0xffaa00 },
  { id:27, name:"Lagos",         country:"Nigeria",      region:"Africa",   lat:6.52,   lng:3.37,    score:68, growth:25, remote:55, salaryIdx:30, jobs:["FinTech Developer","AgriTech Specialist","Mobile Health Engineer"], color:0xff4444 },
  { id:28, name:"Nairobi",       country:"Kenya",        region:"Africa",   lat:-1.29,  lng:36.82,   score:70, growth:22, remote:58, salaryIdx:28, jobs:["Mobile Tech Developer","Climate Adaptation Expert","AgriTech Specialist"], color:0xffaa00 },
  { id:29, name:"Cape Town",     country:"S. Africa",    region:"Africa",   lat:-33.93, lng:18.42,   score:68, growth:14, remote:72, salaryIdx:35, jobs:["Climate Tech Engineer","AgriTech Developer","Remote Work Hub Lead"], color:0xff4444 },
  { id:30, name:"Tashkent",      country:"Uzbekistan",   region:"Central Asia",lat:41.30,lng:69.24,  score:64, growth:28, remote:52, salaryIdx:25, jobs:["IT Outsourcing Lead","FinTech Developer","EdTech Engineer"], color:0xff4444 },
];

const ARCS = [
  [1,19],[1,20],[1,8],[8,11],[8,9],[17,8],[19,23],[2,8],[12,8],[10,13],[1,5],[19,22],[7,2],[28,8],[23,1]
];

const PROFESSIONS = {
  "Software Engineer":     { risk:45, cat:"Tech",        salary:120000, resilient:["system design","problem solving","architecture"], atRisk:["boilerplate coding","code review","debugging","unit tests"] },
  "Data Scientist":        { risk:38, cat:"Tech",        salary:130000, resilient:["statistical strategy","ML architecture","research"], atRisk:["data cleaning","basic analysis","reporting","visualization"] },
  "Cybersecurity Analyst": { risk:28, cat:"Tech",        salary:115000, resilient:["threat hunting","incident response","red teaming"], atRisk:["routine monitoring","log review","basic audits","patching"] },
  "DevOps Engineer":       { risk:42, cat:"Tech",        salary:125000, resilient:["infrastructure strategy","crisis response"], atRisk:["script automation","basic CI/CD","routine deployments"] },
  "UX Designer":           { risk:40, cat:"Creative",    salary:95000,  resilient:["user research","empathy mapping","strategy"], atRisk:["wireframing","basic prototyping","UI generation"] },
  "Graphic Designer":      { risk:75, cat:"Creative",    salary:60000,  resilient:["creative direction","brand strategy","art direction"], atRisk:["image creation","layout design","photo editing","logo design"] },
  "Content Writer":        { risk:80, cat:"Creative",    salary:55000,  resilient:["investigative writing","deep storytelling","journalism"], atRisk:["basic content","SEO articles","social posts","product descriptions"] },
  "Financial Analyst":     { risk:85, cat:"Finance",     salary:90000,  resilient:["strategic interpretation","client relationships"], atRisk:["modeling","Excel tasks","research","forecasting","reporting"] },
  "Accountant":            { risk:94, cat:"Finance",     salary:70000,  resilient:["strategic tax advice","complex audit judgment"], atRisk:["bookkeeping","tax prep","auditing","payroll","data entry"] },
  "Marketing Manager":     { risk:60, cat:"Business",    salary:85000,  resilient:["brand strategy","human insights","creative campaigns"], atRisk:["copywriting","basic analytics","ad targeting","A/B testing"] },
  "Project Manager":       { risk:50, cat:"Business",    salary:95000,  resilient:["leadership","stakeholder management","risk judgment"], atRisk:["scheduling","task tracking","status reports","meeting notes"] },
  "HR Manager":            { risk:55, cat:"Business",    salary:80000,  resilient:["culture building","conflict resolution","leadership"], atRisk:["resume screening","compliance checks","onboarding","payroll"] },
  "Doctor":                { risk:25, cat:"Healthcare",  salary:200000, resilient:["surgery","clinical judgment","empathy","diagnosis"], atRisk:["imaging analysis","routine diagnosis","documentation"] },
  "Nurse":                 { risk:20, cat:"Healthcare",  salary:75000,  resilient:["patient care","empathy","emergency response"], atRisk:["documentation","basic monitoring","scheduling"] },
  "Psychologist":          { risk:15, cat:"Healthcare",  salary:90000,  resilient:["therapy","empathy","human connection","trauma care"], atRisk:["basic assessments","intake forms","routine CBT"] },
  "Teacher":               { risk:30, cat:"Education",   salary:55000,  resilient:["mentoring","emotional support","inspiration","leadership"], atRisk:["content delivery","basic grading","quiz creation"] },
  "Lawyer":                { risk:55, cat:"Legal",       salary:140000, resilient:["negotiation","case strategy","courtroom advocacy"], atRisk:["legal research","contract drafting","document review"] },
  "Truck Driver":          { risk:90, cat:"Transport",   salary:50000,  resilient:["emergency judgment","complex urban navigation"], atRisk:["long-haul driving","route planning","navigation","scheduling"] },
  "Electrician":           { risk:18, cat:"Trades",      salary:65000,  resilient:["complex installation","safety judgment","troubleshooting"], atRisk:["basic diagnostics","permit documentation"] },
  "Chef":                  { risk:22, cat:"Hospitality", salary:55000,  resilient:["culinary creativity","guest experience","innovation"], atRisk:["prep work","standard recipes","ordering","scheduling"] },
};

const PIVOTS = {
  "Software Engineer":     [{to:"AI Ethics Officer",match:85,emoji:"⚖️",months:6},{to:"Quantum Developer",match:78,emoji:"⚛️",months:12},{to:"Neural Interface Developer",match:74,emoji:"🧠",months:18},{to:"Prompt Engineer",match:90,emoji:"💬",months:3},{to:"Climate Tech Engineer",match:68,emoji:"🌿",months:12}],
  "Data Scientist":        [{to:"Longevity Scientist",match:88,emoji:"🧬",months:8},{to:"Climate Tech Engineer",match:84,emoji:"🌿",months:6},{to:"Biosecurity Analyst",match:80,emoji:"🔬",months:10},{to:"AI Ethics Officer",match:76,emoji:"⚖️",months:6},{to:"Quantum Developer",match:70,emoji:"⚛️",months:14}],
  "Lawyer":                [{to:"AI Ethics Officer",match:92,emoji:"⚖️",months:4},{to:"Biosecurity Analyst",match:78,emoji:"🔬",months:10},{to:"Human-AI Collaboration Spec.",match:72,emoji:"🤝",months:6},{to:"Green Energy Consultant",match:70,emoji:"☀️",months:8},{to:"Digital Wellness Coach",match:62,emoji:"🧘",months:5}],
  "Doctor":                [{to:"Longevity Scientist",match:94,emoji:"🧬",months:12},{to:"Biosecurity Analyst",match:88,emoji:"🔬",months:8},{to:"AI Ethics Officer",match:72,emoji:"⚖️",months:6},{to:"Digital Wellness Coach",match:75,emoji:"🧘",months:4},{to:"Neural Interface Developer",match:68,emoji:"🧠",months:18}],
  "Teacher":               [{to:"Human-AI Collaboration Spec.",match:86,emoji:"🤝",months:5},{to:"Digital Wellness Coach",match:82,emoji:"🧘",months:4},{to:"Prompt Engineer",match:72,emoji:"💬",months:3},{to:"AI Ethics Officer",match:70,emoji:"⚖️",months:8},{to:"Metaverse Architect",match:65,emoji:"🌐",months:14}],
  "Financial Analyst":     [{to:"AI Ethics Officer",match:80,emoji:"⚖️",months:6},{to:"Green Energy Consultant",match:75,emoji:"☀️",months:8},{to:"Prompt Engineer",match:68,emoji:"💬",months:3},{to:"Human-AI Collaboration Spec.",match:72,emoji:"🤝",months:5},{to:"FinTech Architect",match:82,emoji:"💳",months:7}],
};
const getPivots = p => PIVOTS[p] || Object.values(PIVOTS)[0];

const FUTURE_JOBS = [
  { title:"AI Ethics Officer",             emoji:"⚖️", demand:95, growth:"+340%", salary:"$140–200k", desc:"Govern AI systems for fairness, transparency and legal compliance" },
  { title:"Prompt Engineer",               emoji:"💬", demand:92, growth:"+290%", salary:"$90–160k",  desc:"Design and optimize AI systems across industries and domains" },
  { title:"Climate Tech Engineer",         emoji:"🌿", demand:90, growth:"+280%", salary:"$110–160k", desc:"Engineer large-scale climate change mitigation solutions" },
  { title:"Human-AI Collaboration Spec.",  emoji:"🤝", demand:89, growth:"+260%", salary:"$100–150k", desc:"Bridge human workers and intelligent AI systems in organizations" },
  { title:"Quantum Developer",             emoji:"⚛️", demand:85, growth:"+310%", salary:"$150–220k", desc:"Build next-gen applications on quantum computing platforms" },
  { title:"Longevity Scientist",           emoji:"🧬", demand:80, growth:"+220%", salary:"$130–190k", desc:"Research anti-aging therapies and life extension biotechnologies" },
  { title:"Biosecurity Analyst",           emoji:"🔬", demand:80, growth:"+200%", salary:"$110–160k", desc:"Protect nations and organizations from emerging biological threats" },
  { title:"Green Energy Consultant",       emoji:"☀️", demand:88, growth:"+240%", salary:"$95–145k",  desc:"Guide organizations through full renewable energy transitions" },
  { title:"Neural Interface Developer",    emoji:"🧠", demand:72, growth:"+380%", salary:"$160–240k", desc:"Build software for brain-computer interface platforms" },
  { title:"Digital Wellness Coach",        emoji:"🧘", demand:76, growth:"+175%", salary:"$70–110k",  desc:"Help individuals maintain wellbeing in a hyper-digital world" },
  { title:"Metaverse Architect",           emoji:"🌐", demand:70, growth:"+190%", salary:"$120–180k", desc:"Design immersive virtual worlds and spatial computing experiences" },
  { title:"Space Tourism Specialist",      emoji:"🚀", demand:65, growth:"+450%", salary:"$90–140k",  desc:"Design and manage commercial space travel experiences" },
];

// ── HELPERS ──────────────────────────────────────────────────────────────────
function ll2v3(lat, lng, r=1.015) {
  const phi   = (90-lat)*(Math.PI/180);
  const theta = (lng+180)*(Math.PI/180);
  return new THREE.Vector3(-r*Math.sin(phi)*Math.cos(theta), r*Math.cos(phi), r*Math.sin(phi)*Math.sin(theta));
}
const scoreColor = s => s>=90?"#00ff88":s>=80?"#00aaff":s>=70?"#ffaa00":"#ff4444";
const riskColor  = r => r<30?"#00ff88":r<55?"#ffcc00":r<75?"#ff8844":"#ff3344";
const riskLabel  = r => r<30?"🟢 Low Risk":r<55?"🟡 Medium Risk":r<75?"🟠 High Risk":"🔴 Critical Risk";

// ── GLOBE SCENE ──────────────────────────────────────────────────────────────
function buildScene(el, onCityClick) {
  const W=el.clientWidth, H=el.clientHeight;
  const scene    = new THREE.Scene();
  const camera   = new THREE.PerspectiveCamera(42, W/H, 0.1, 500);
  camera.position.z = 2.9;

  const renderer = new THREE.WebGLRenderer({antialias:true, alpha:true});
  renderer.setSize(W, H);
  renderer.setPixelRatio(Math.min(devicePixelRatio, 2));
  el.appendChild(renderer.domElement);

  // Stars
  const sg = new THREE.BufferGeometry();
  sg.setAttribute("position", new THREE.BufferAttribute(new Float32Array(4000*3).map(()=>(Math.random()-.5)*180),3));
  scene.add(new THREE.Points(sg, new THREE.PointsMaterial({color:0xffffff,size:0.06,transparent:true,opacity:0.55})));

  // Globe group
  const globe = new THREE.Group();
  scene.add(globe);

  // Core sphere — gradient-like dark blue
  const coreMat = new THREE.MeshPhongMaterial({color:0x051525, emissive:0x0a1c35, shininess:60, specular:0x0044aa});
  globe.add(new THREE.Mesh(new THREE.SphereGeometry(1,64,64), coreMat));

  // Tectonic-style wireframe
  globe.add(Object.assign(new THREE.Mesh(new THREE.SphereGeometry(1.0015,30,30),
    new THREE.MeshBasicMaterial({color:0x0a3a6a, wireframe:true, opacity:0.08, transparent:true})),{renderOrder:1}));

  // Atmosphere glow (outer)
  globe.add(new THREE.Mesh(new THREE.SphereGeometry(1.1,32,32),
    new THREE.MeshPhongMaterial({color:0x0055cc, transparent:true, opacity:0.055, side:THREE.BackSide})));

  // Atmosphere glow (inner rim)
  globe.add(new THREE.Mesh(new THREE.SphereGeometry(1.04,32,32),
    new THREE.MeshPhongMaterial({color:0x003399, transparent:true, opacity:0.04, side:THREE.BackSide})));

  // Lights
  scene.add(new THREE.AmbientLight(0x1a3355, 1.3));
  const sun = new THREE.DirectionalLight(0x5599ff, 2.5);
  sun.position.set(5,3,5); scene.add(sun);
  const fill = new THREE.DirectionalLight(0x002244, 0.4);
  fill.position.set(-4,-2,-3); scene.add(fill);

  // City markers + glows
  const markerMeshes = [];
  CITIES.forEach(c => {
    const pos = ll2v3(c.lat, c.lng, 1.018);
    // outer glow ring
    const glow = new THREE.Mesh(
      new THREE.SphereGeometry(0.022, 8, 8),
      new THREE.MeshBasicMaterial({color:c.color, transparent:true, opacity:0.3})
    );
    glow.position.copy(pos);
    glow.userData = {isGlow:true};
    globe.add(glow);
    // core dot
    const dot = new THREE.Mesh(
      new THREE.SphereGeometry(0.011, 8, 8),
      new THREE.MeshBasicMaterial({color:0xffffff})
    );
    dot.position.copy(pos);
    dot.userData = {...c, isMarker:true};
    globe.add(dot);
    markerMeshes.push(dot, glow);
  });

  // Arc curves between cities
  const arcGroup = new THREE.Group();
  globe.add(arcGroup);
  ARCS.forEach(([aId, bId]) => {
    const a = CITIES.find(c=>c.id===aId), b = CITIES.find(c=>c.id===bId);
    if (!a||!b) return;
    const va = ll2v3(a.lat, a.lng, 1), vb = ll2v3(b.lat, b.lng, 1);
    const mid = va.clone().add(vb).normalize().multiplyScalar(1.28);
    const curve = new THREE.QuadraticBezierCurve3(va, mid, vb);
    const pts = curve.getPoints(60);
    const geo = new THREE.BufferGeometry().setFromPoints(pts);
    const mat = new THREE.LineBasicMaterial({color:0x005588, transparent:true, opacity:0.3});
    arcGroup.add(new THREE.Line(geo, mat));
  });

  // Pointer
  let dragging=false, lastX=0, lastY=0, downX=0, downY=0;
  const t = {v:0};

  const onMD = e => { dragging=true; lastX=downX=e.clientX; lastY=downY=e.clientY; el.style.cursor="grabbing"; };
  const onMM = e => {
    if (!dragging) return;
    globe.rotation.y += (e.clientX-lastX)*0.005;
    globe.rotation.x = Math.max(-1.2, Math.min(1.2, globe.rotation.x+(e.clientY-lastY)*0.005));
    lastX=e.clientX; lastY=e.clientY;
  };
  const onMU = e => {
    el.style.cursor="grab"; dragging=false;
    if (Math.hypot(e.clientX-downX, e.clientY-downY)<5) {
      const rect = el.getBoundingClientRect();
      const mouse = new THREE.Vector2(((e.clientX-rect.left)/rect.width)*2-1, -((e.clientY-rect.top)/rect.height)*2+1);
      const rc = new THREE.Raycaster(); rc.setFromCamera(mouse, camera);
      const hits = rc.intersectObjects(markerMeshes.filter(m=>m.userData.isMarker));
      if (hits.length) onCityClick(hits[0].object.userData);
    }
  };
  const onWheel = e => {
    camera.position.z = Math.max(1.8, Math.min(4.5, camera.position.z + e.deltaY*0.003));
  };
  el.addEventListener("mousedown", onMD);
  el.addEventListener("mousemove", onMM);
  el.addEventListener("mouseup",   onMU);
  el.addEventListener("wheel",     onWheel, {passive:true});

  const onResize = () => {
    const W=el.clientWidth, H=el.clientHeight;
    camera.aspect=W/H; camera.updateProjectionMatrix();
    renderer.setSize(W,H);
  };
  window.addEventListener("resize", onResize);

  let raf, tick=0;
  const animate = () => {
    raf = requestAnimationFrame(animate);
    tick += 0.015;
    if (!dragging) globe.rotation.y += 0.0006;
    markerMeshes.forEach(m => {
      if (m.userData.isGlow) m.material.opacity = 0.15 + 0.2*Math.abs(Math.sin(tick+m.id||0));
    });
    arcGroup.children.forEach((l,i) => {
      l.material.opacity = 0.12 + 0.22*Math.abs(Math.sin(tick*0.8+i));
    });
    renderer.render(scene, camera);
  };
  animate();

  return () => {
    cancelAnimationFrame(raf);
    window.removeEventListener("resize", onResize);
    el.removeEventListener("mousedown",onMD);
    el.removeEventListener("mousemove",onMM);
    el.removeEventListener("mouseup",onMU);
    el.removeEventListener("wheel",onWheel);
    if (el.contains(renderer.domElement)) el.removeChild(renderer.domElement);
    renderer.dispose();
  };
}

// ── COMPONENTS ───────────────────────────────────────────────────────────────
function Bar({value, color="#00aaff", h=6}) {
  return (
    <div style={{height:h, background:"#0d2137", borderRadius:h, overflow:"hidden"}}>
      <div style={{width:`${value}%`, height:"100%", background:color, borderRadius:h, transition:"width 0.6s ease"}}/>
    </div>
  );
}

function Tag({label}) {
  return <span style={{padding:"3px 10px", background:"#0d2040", borderRadius:20, fontSize:11, color:"#2a6a9a", border:"1px solid #0d2a40"}}>{label}</span>;
}

function StatBox({icon, value, label, color="#00ff88"}) {
  return (
    <div style={{flex:1, padding:"12px 8px", background:"#071525", borderRadius:12, textAlign:"center", border:"1px solid #0d2137"}}>
      <div style={{fontSize:16, marginBottom:4}}>{icon}</div>
      <div style={{fontSize:18, fontWeight:900, color}}>{value}</div>
      <div style={{fontSize:10, color:"#1a4060", marginTop:2, letterSpacing:.5}}>{label}</div>
    </div>
  );
}

// ── CITY PANEL ───────────────────────────────────────────────────────────────
function CityPanel({city, onClose}) {
  const sc = scoreColor(city.score);
  return (
    <div style={{position:"absolute", right:16, top:16, width:320, background:"rgba(3,10,22,0.97)", border:`1px solid ${sc}33`, borderRadius:20, padding:22, backdropFilter:"blur(18px)", boxShadow:`0 0 40px ${sc}18`}}>
      <div style={{display:"flex", justifyContent:"space-between", alignItems:"flex-start", marginBottom:18}}>
        <div>
          <div style={{fontSize:22, fontWeight:900, lineHeight:1.1, color:"#eef"}}>{city.name}</div>
          <div style={{fontSize:12, color:"#1a4a6a", marginTop:3}}>{city.country} · {city.region}</div>
        </div>
        <button onClick={onClose} style={{background:"none",border:"none",color:"#1a4a6a",cursor:"pointer",fontSize:24,lineHeight:1,padding:0,marginTop:-2}}>×</button>
      </div>

      <div style={{marginBottom:16}}>
        <div style={{display:"flex", justifyContent:"space-between", marginBottom:5}}>
          <span style={{fontSize:11, color:"#1a4060", letterSpacing:1}}>OPPORTUNITY SCORE</span>
          <span style={{fontSize:18, fontWeight:900, color:sc}}>{city.score}/100</span>
        </div>
        <Bar value={city.score} color={`linear-gradient(to right, ${sc}99, ${sc})`} h={8}/>
      </div>

      <div style={{display:"flex", gap:8, marginBottom:18}}>
        <StatBox icon="📈" value={`+${city.growth}%`} label="JOB GROWTH" color="#00ff88"/>
        <StatBox icon="🌐" value={`${city.remote}%`}  label="REMOTE" color="#00aaff"/>
        <StatBox icon="💰" value={`${city.salaryIdx}`} label="SALARY IDX" color="#ffaa00"/>
      </div>

      <div style={{marginBottom:14}}>
        <div style={{fontSize:10, color:"#1a4060", marginBottom:9, letterSpacing:1.2}}>🎯 TOP FUTURE CAREERS</div>
        {city.jobs.map((j,i)=>(
          <div key={i} style={{padding:"8px 12px", marginBottom:5, background:"#071525", borderRadius:10, fontSize:13, display:"flex", gap:10, alignItems:"center", border:"1px solid #0d2137"}}>
            <span style={{fontSize:16}}>{["🥇","🥈","🥉"][i]}</span>
            <span style={{color:"#8ab4d4", fontWeight:i===0?700:400}}>{j}</span>
          </div>
        ))}
      </div>

      <div style={{fontSize:10, color:"#0d2a3a", textAlign:"center", paddingTop:8, borderTop:"1px solid #0d2030"}}>
        Data refreshed daily · AI-powered insights
      </div>
    </div>
  );
}

// ── CAREER ANALYZER ───────────────────────────────────────────────────────────
function CareerAnalyzer() {
  const [prof, setProf] = useState("");
  const [result, setResult] = useState(null);
  const [loading, setLoading] = useState(false);
  const analyze = () => {
    if (!PROFESSIONS[prof]) return;
    setLoading(true);
    setTimeout(() => {
      setResult({prof, ...PROFESSIONS[prof], pivots: getPivots(prof)});
      setLoading(false);
    }, 900);
  };
  const rc = result ? riskColor(result.risk) : "#00ff88";

  return (
    <div style={{maxWidth:860, margin:"0 auto"}}>
      <div style={{marginBottom:32}}>
        <h2 style={{fontSize:28, fontWeight:900, margin:"0 0 8px", background:"linear-gradient(to right,#00aaff,#00ff88)", WebkitBackgroundClip:"text", WebkitTextFillColor:"transparent"}}>
          AI Replacement Risk Analyzer
        </h2>
        <p style={{color:"#1a4a6a", margin:0, fontSize:15}}>Discover your displacement risk and the best career pivots available to you</p>
      </div>

      <div style={{display:"flex", gap:12, marginBottom:32}}>
        <select value={prof} onChange={e=>setProf(e.target.value)} style={{flex:1, padding:"14px 18px", background:"#071525", border:"1px solid #0d2a47", borderRadius:14, color:prof?"#cde":"#1a4a6a", fontSize:15, cursor:"pointer", outline:"none"}}>
          <option value="">— Select your current profession —</option>
          {Object.entries(PROFESSIONS).map(([p,d])=>(<option key={p} value={p}>{p} — {d.cat}</option>))}
        </select>
        <button onClick={analyze} disabled={!prof||loading} style={{padding:"14px 30px", borderRadius:14, border:"none", cursor:prof&&!loading?"pointer":"not-allowed", fontSize:15, fontWeight:800, background:prof?"linear-gradient(135deg,#00aaff,#00ff88)":"#0d2137", color:prof?"#020b18":"#2a5a7a", minWidth:140, transition:"all 0.2s"}}>
          {loading?"Analyzing…":"Analyze →"}
        </button>
      </div>

      {result && (
        <div style={{animation:"fadeIn 0.4s ease"}}>
          {/* Risk overview */}
          <div style={{padding:28, background:"#071525", border:`2px solid ${rc}28`, borderRadius:20, marginBottom:24, boxShadow:`0 0 40px ${rc}0a`}}>
            <div style={{display:"flex", justifyContent:"space-between", alignItems:"flex-start", marginBottom:20, gap:20, flexWrap:"wrap"}}>
              <div style={{flex:1, minWidth:240}}>
                <div style={{fontSize:26, fontWeight:900, color:"#eef"}}>{result.prof}</div>
                <div style={{color:"#1a4a6a", marginTop:4, fontSize:13}}>{result.cat} · Avg ${result.salary.toLocaleString()}/yr</div>
                <div style={{marginTop:14, fontSize:14, color:"#5a8aaa", lineHeight:1.7, maxWidth:430}}>
                  {result.risk<30?"Strong AI resilience. Your human-centric skills create a wide moat. Deepen them strategically.":
                   result.risk<55?"Moderate risk zone. Core skills stay valuable but automation is encroaching. Upskill now.":
                   result.risk<75?"High risk. Significant portions of your role will be automated within 3–5 years. Act soon.":
                   "Critical risk. Most of this role will be automated in 2–4 years. An immediate pivot is strongly advised."}
                </div>
              </div>
              <div style={{textAlign:"center", flexShrink:0}}>
                <div style={{width:110, height:110, borderRadius:"50%", background:`conic-gradient(${rc} ${result.risk*3.6}deg, #0d2137 0)`, display:"flex", alignItems:"center", justifyContent:"center", margin:"0 auto"}}>
                  <div style={{width:84, height:84, borderRadius:"50%", background:"#071525", display:"flex", flexDirection:"column", alignItems:"center", justifyContent:"center"}}>
                    <span style={{fontSize:26, fontWeight:900, color:rc, lineHeight:1}}>{result.risk}%</span>
                    <span style={{fontSize:9, color:"#1a4060", letterSpacing:.5}}>RISK</span>
                  </div>
                </div>
                <div style={{fontSize:12, marginTop:8, color:rc}}>{riskLabel(result.risk)}</div>
              </div>
            </div>

            <div style={{display:"grid", gridTemplateColumns:"1fr 1fr", gap:14}}>
              {[["✅ RESILIENT SKILLS","#00ff88", result.resilient],["⚠️ AT-RISK TASKS","#ff5544", result.atRisk]].map(([lbl,col,items])=>(
                <div key={lbl} style={{padding:16, background:"#0a1c2e", borderRadius:14, border:`1px solid ${col}22`}}>
                  <div style={{fontSize:10, color:col, marginBottom:10, letterSpacing:1.2, fontWeight:700}}>{lbl}</div>
                  {items.map(s=><div key={s} style={{fontSize:13, color:"#6a9aba", marginBottom:6, display:"flex", gap:7, alignItems:"center"}}>
                    <div style={{width:5, height:5, borderRadius:"50%", background:col, flexShrink:0}}/>
                    {s}
                  </div>)}
                </div>
              ))}
            </div>
          </div>

          {/* Pivots */}
          <h3 style={{fontSize:20, fontWeight:800, margin:"0 0 16px", color:"#00aaff"}}>🔄 Your Best Career Pivots</h3>
          {result.pivots.map((t,i)=>(
            <div key={i} style={{padding:18, background:"#071525", border:"1px solid #0d2137", borderRadius:16, marginBottom:12, display:"flex", gap:16, alignItems:"center", transition:"all 0.2s", cursor:"default"}}
              onMouseEnter={e=>{e.currentTarget.style.borderColor="#00aaff44"; e.currentTarget.style.transform="translateX(4px)";}}
              onMouseLeave={e=>{e.currentTarget.style.borderColor="#0d2137"; e.currentTarget.style.transform="translateX(0)";}}>
              <div style={{width:60, height:60, borderRadius:"50%", flexShrink:0, background:`conic-gradient(#00ff88 ${t.match*3.6}deg,#0d2137 0)`, display:"flex", alignItems:"center", justifyContent:"center"}}>
                <div style={{width:46, height:46, borderRadius:"50%", background:"#071525", display:"flex", alignItems:"center", justifyContent:"center"}}>
                  <span style={{fontSize:13, fontWeight:900, color:"#00ff88"}}>{t.match}%</span>
                </div>
              </div>
              <div style={{flex:1}}>
                <div style={{display:"flex", justifyContent:"space-between", alignItems:"center", marginBottom:4}}>
                  <span style={{fontSize:16, fontWeight:700}}>{t.emoji} {t.to}</span>
                  <span style={{fontSize:11, color:"#1a4a6a", flexShrink:0, marginLeft:8}}>~{t.months}mo reskill</span>
                </div>
                <div style={{fontSize:13, color:"#2a5a7a"}}>Skill match: strong transferability from your {result.cat} background</div>
              </div>
            </div>
          ))}
        </div>
      )}
      <style>{`@keyframes fadeIn{from{opacity:0;transform:translateY(12px)}to{opacity:1;transform:translateY(0)}}`}</style>
    </div>
  );
}

// ── FUTURE JOBS ───────────────────────────────────────────────────────────────
function FutureJobs() {
  const [filter, setFilter] = useState("All");
  const cats = ["All","Tech","Life Science","Green","Wellness","Space","Creative"];
  const catMap = {
    "AI Ethics Officer":"Tech","Prompt Engineer":"Tech","Quantum Developer":"Tech","Neural Interface Developer":"Tech","Metaverse Architect":"Creative",
    "Climate Tech Engineer":"Green","Green Energy Consultant":"Green",
    "Longevity Scientist":"Life Science","Biosecurity Analyst":"Life Science",
    "Digital Wellness Coach":"Wellness","Human-AI Collaboration Spec.":"Wellness",
    "Space Tourism Specialist":"Space",
  };
  const filtered = filter==="All"?FUTURE_JOBS:FUTURE_JOBS.filter(j=>catMap[j.title]===filter);
  return (
    <div style={{maxWidth:1100, margin:"0 auto"}}>
      <div style={{marginBottom:28}}>
        <h2 style={{fontSize:28, fontWeight:900, margin:"0 0 8px", background:"linear-gradient(to right,#00aaff,#00ff88)", WebkitBackgroundClip:"text", WebkitTextFillColor:"transparent"}}>
          Top Future Professions 2025–2035
        </h2>
        <p style={{color:"#1a4a6a", margin:"0 0 20px", fontSize:15}}>The most demanded careers in the post-AI world</p>
        <div style={{display:"flex", gap:8, flexWrap:"wrap"}}>
          {cats.map(c=>(
            <button key={c} onClick={()=>setFilter(c)} style={{padding:"6px 16px", borderRadius:20, border:`1px solid ${filter===c?"#00aaff":"#0d2030"}`, background:filter===c?"rgba(0,170,255,0.12)":"transparent", color:filter===c?"#00aaff":"#2a5a7a", cursor:"pointer", fontSize:12, fontWeight:filter===c?700:400, transition:"all 0.15s"}}>{c}</button>
          ))}
        </div>
      </div>
      <div style={{display:"grid", gridTemplateColumns:"repeat(auto-fill,minmax(290px,1fr))", gap:16}}>
        {filtered.map((j,i)=>(
          <div key={i} style={{padding:22, background:"#071525", border:"1px solid #0d2137", borderRadius:18, transition:"all 0.2s", cursor:"default"}}
            onMouseEnter={e=>{e.currentTarget.style.borderColor="#00aaff44"; e.currentTarget.style.transform="translateY(-4px)"; e.currentTarget.style.boxShadow="0 12px 40px #00aaff12";}}
            onMouseLeave={e=>{e.currentTarget.style.borderColor="#0d2137"; e.currentTarget.style.transform="translateY(0)"; e.currentTarget.style.boxShadow="none";}}>
            <div style={{display:"flex", justifyContent:"space-between", alignItems:"flex-start", marginBottom:10}}>
              <div style={{display:"flex", gap:10, alignItems:"center"}}>
                <span style={{fontSize:28}}>{j.emoji}</span>
                <div style={{fontSize:15, fontWeight:700, lineHeight:1.3, color:"#cde"}}>{j.title}</div>
              </div>
              <div style={{fontSize:13, color:"#00ff88", fontWeight:900, flexShrink:0, marginLeft:8}}>{j.growth}</div>
            </div>
            <div style={{fontSize:13, color:"#2a5a7a", marginBottom:14, lineHeight:1.6}}>{j.desc}</div>
            <div style={{marginBottom:12}}>
              <div style={{display:"flex", justifyContent:"space-between", marginBottom:5}}>
                <span style={{fontSize:10, color:"#1a4060", letterSpacing:1}}>DEMAND INDEX</span>
                <span style={{fontSize:11, color:"#00aaff", fontWeight:700}}>{j.demand}/100</span>
              </div>
              <Bar value={j.demand} color="linear-gradient(to right,#00aaff,#00ff88)"/>
            </div>
            <div style={{fontSize:12, color:"#1a5a7a", background:"#0a1c2e", padding:"6px 12px", borderRadius:8, textAlign:"center"}}>
              💰 {j.salary}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

// ── GLOBE TAB ──────────────────────────────────────────────────────────────
function GlobeTab() {
  const mountRef  = useRef(null);
  const cleanupRef = useRef(null);
  const [city, setCity]     = useState(null);
  const [search, setSearch] = useState("");
  const [region, setRegion] = useState("All");
  const regions = ["All", ...new Set(CITIES.map(c=>c.region))];

  const onCityClick = useCallback(c => setCity(c), []);

  useEffect(() => {
    if (!mountRef.current) return;
    cleanupRef.current = buildScene(mountRef.current, onCityClick);
    return () => cleanupRef.current?.();
  }, []);

  const filtered = CITIES.filter(c =>
    (region==="All"||c.region===region) &&
    c.name.toLowerCase().includes(search.toLowerCase())
  ).sort((a,b)=>b.score-a.score);

  return (
    <div style={{flex:1, position:"relative", overflow:"hidden"}}>
      <div ref={mountRef} style={{width:"100%", height:"100%", cursor:"grab"}}/>

      {/* Search + Filter */}
      <div style={{position:"absolute", top:16, left:16, display:"flex", flexDirection:"column", gap:8, width:230}}>
        <input value={search} onChange={e=>setSearch(e.target.value)} placeholder="🔍 Search city…"
          style={{padding:"9px 14px", background:"rgba(3,10,22,0.9)", border:"1px solid #0d2a3a", borderRadius:12, color:"#cde", fontSize:13, outline:"none", backdropFilter:"blur(12px)"}}/>
        <select value={region} onChange={e=>setRegion(e.target.value)}
          style={{padding:"8px 14px", background:"rgba(3,10,22,0.9)", border:"1px solid #0d2a3a", borderRadius:12, color:"#8ab4d4", fontSize:12, outline:"none", cursor:"pointer", backdropFilter:"blur(12px)"}}>
          {regions.map(r=><option key={r} value={r}>{r}</option>)}
        </select>
      </div>

      {/* Legend */}
      <div style={{position:"absolute", bottom:16, left:16, background:"rgba(3,10,22,0.9)", border:"1px solid #0c2035", borderRadius:14, padding:"12px 16px", backdropFilter:"blur(12px)"}}>
        <div style={{fontSize:9, color:"#1a4060", marginBottom:8, letterSpacing:2}}>OPPORTUNITY INDEX</div>
        {[["#00ff88","90–100","Elite Hub"],["#00aaff","80–89","High Growth"],["#ffaa00","70–79","Emerging"],["#ff4444","< 70","Developing"]].map(([c,r,l])=>(
          <div key={r} style={{display:"flex", alignItems:"center", gap:8, marginBottom:4, fontSize:12}}>
            <div style={{width:8, height:8, borderRadius:"50%", background:c, flexShrink:0, boxShadow:`0 0 6px ${c}`}}/>
            <span style={{color:"#8ab4d4", width:44, fontSize:11}}>{r}</span>
            <span style={{color:"#1a4060", fontSize:11}}>{l}</span>
          </div>
        ))}
      </div>

      {/* City list */}
      <div style={{position:"absolute", bottom:16, right:city?354:16, top:16, width:220, background:"rgba(3,10,22,0.88)", border:"1px solid #0c2035", borderRadius:16, padding:"14px 10px", backdropFilter:"blur(12px)", overflowY:"auto", transition:"right 0.3s ease"}}>
        <div style={{fontSize:9, color:"#1a4060", marginBottom:10, letterSpacing:2, paddingLeft:4}}>
          🏆 {region==="All"?"GLOBAL":"REGION"} RANKINGS ({filtered.length})
        </div>
        {filtered.map((c,i)=>(
          <div key={c.id} onClick={()=>setCity(c)} style={{display:"flex", alignItems:"center", gap:8, padding:"7px 8px", marginBottom:3, borderRadius:10, cursor:"pointer", transition:"background 0.15s", border:"1px solid transparent"}}
            onMouseEnter={e=>{e.currentTarget.style.background="rgba(0,170,255,0.08)"; e.currentTarget.style.borderColor="#00aaff22";}}
            onMouseLeave={e=>{e.currentTarget.style.background="transparent"; e.currentTarget.style.borderColor="transparent";}}>
            <span style={{fontSize:10, color:"#1a4060", width:22, textAlign:"center"}}>#{i+1}</span>
            <div style={{flex:1, minWidth:0}}>
              <div style={{fontSize:12, color:"#cde", fontWeight:city?.id===c.id?700:400, whiteSpace:"nowrap", overflow:"hidden", textOverflow:"ellipsis"}}>{c.name}</div>
              <div style={{fontSize:10, color:"#1a4060"}}>{c.country}</div>
            </div>
            <span style={{fontSize:12, fontWeight:800, color:scoreColor(c.score), flexShrink:0}}>{c.score}</span>
          </div>
        ))}
      </div>

      {city && <CityPanel city={city} onClose={()=>setCity(null)}/>}

      <div style={{position:"absolute", bottom:16, left:"50%", transform:"translateX(-50%)", fontSize:11, color:"#0d2233", pointerEvents:"none", whiteSpace:"nowrap"}}>
        Drag to rotate · Scroll to zoom · Click markers or list to explore
      </div>
    </div>
  );
}

// ── MAIN APP ──────────────────────────────────────────────────────────────────
export default function App() {
  const [tab, setTab] = useState("globe");

  const NAV = [
    {id:"globe",  label:"🌍 Globe"},
    {id:"career", label:"⚡ Career Analyzer"},
    {id:"future", label:"🚀 Future Jobs"},
  ];

  return (
    <div style={{minHeight:"100vh", background:"#030a12", color:"#fff", fontFamily:"'Inter','Segoe UI',system-ui,sans-serif", display:"flex", flexDirection:"column"}}>

      {/* Header */}
      <header style={{padding:"0 24px", height:60, display:"flex", alignItems:"center", justifyContent:"space-between", borderBottom:"1px solid #0c1e30", background:"rgba(3,10,18,0.98)", backdropFilter:"blur(16px)", position:"sticky", top:0, zIndex:100, flexShrink:0}}>
        <div style={{display:"flex", alignItems:"center", gap:12}}>
          <div style={{width:36, height:36, borderRadius:"50%", background:"linear-gradient(135deg,#0055aa,#00ff88)", display:"flex", alignItems:"center", justifyContent:"center", fontSize:18, flexShrink:0}}>🌍</div>
          <div>
            <div style={{fontSize:17, fontWeight:900, background:"linear-gradient(to right,#00aaff,#00ff88)", WebkitBackgroundClip:"text", WebkitTextFillColor:"transparent", lineHeight:1.1}}>FutureProof</div>
            <div style={{fontSize:8, color:"#0d3a5a", letterSpacing:2.5, lineHeight:1}}>GLOBAL CAREER INTELLIGENCE</div>
          </div>
        </div>
        <nav style={{display:"flex", gap:4}}>
          {NAV.map(({id,label})=>(
            <button key={id} onClick={()=>setTab(id)} style={{padding:"7px 18px", borderRadius:20, border:`1px solid ${tab===id?"#00aaff55":"transparent"}`, background:tab===id?"rgba(0,170,255,0.1)":"transparent", color:tab===id?"#00aaff":"#2a5a7a", cursor:"pointer", fontSize:13, fontWeight:tab===id?700:400, transition:"all 0.2s"}}>
              {label}
            </button>
          ))}
        </nav>
        <div style={{display:"flex", alignItems:"center", gap:8}}>
          <div style={{width:7, height:7, borderRadius:"50%", background:"#00ff88", boxShadow:"0 0 8px #00ff88", animation:"pulse 2s infinite"}}/>
          <span style={{fontSize:11, color:"#1a4060"}}>Live Data</span>
        </div>
      </header>

      {/* Body */}
      <div style={{flex:1, display:"flex", overflow:"hidden", height:"calc(100vh - 60px)"}}>
        {tab==="globe"  && <GlobeTab/>}
        {tab==="career" && <div style={{flex:1, padding:"36px 40px", overflowY:"auto"}}><CareerAnalyzer/></div>}
        {tab==="future" && <div style={{flex:1, padding:"36px 40px", overflowY:"auto"}}><FutureJobs/></div>}
      </div>

      <style>{`
        @keyframes pulse { 0%,100%{opacity:1} 50%{opacity:0.4} }
        ::-webkit-scrollbar { width:4px; background:#030a12 }
        ::-webkit-scrollbar-thumb { background:#0d2a3a; border-radius:4px }
        select option { background:#071525; color:#cde }
      `}</style>
    </div>
  );
}
