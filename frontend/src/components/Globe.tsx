import { useRef, useMemo, useEffect } from 'react';
import { Canvas, useFrame } from '@react-three/fiber';
import { OrbitControls, Sphere } from '@react-three/drei';
import * as THREE from 'three';
import { useGlobeStore } from '../store';
import type { CityWithScore } from '../types';

// City marker component
function CityMarker({ city, onClick }: { city: CityWithScore; onClick: (city: CityWithScore) => void }) {
  const meshRef = useRef<THREE.Mesh>(null);

  // Convert lat/lng to 3D coordinates
  const position = useMemo(() => {
    const phi = (90 - city.lat) * (Math.PI / 180);
    const theta = (city.lng + 180) * (Math.PI / 180);
    const radius = 2.02;

    return new THREE.Vector3(
      -radius * Math.sin(phi) * Math.cos(theta),
      radius * Math.cos(phi),
      radius * Math.sin(phi) * Math.sin(theta)
    );
  }, [city.lat, city.lng]);

  // Color based on score
  const color = useMemo(() => {
    if (city.score >= 85) return '#10b981'; // Green
    if (city.score >= 70) return '#f59e0b'; // Yellow
    return '#ef4444'; // Red
  }, [city.score]);

  return (
    <mesh
      ref={meshRef}
      position={position}
      onClick={() => onClick(city)}
    >
      <sphereGeometry args={[0.02, 16, 16]} />
      <meshBasicMaterial color={color} />
    </mesh>
  );
}

// Rotating Earth
function Earth() {
  const meshRef = useRef<THREE.Mesh>(null);

  useFrame(() => {
    if (meshRef.current) {
      meshRef.current.rotation.y += 0.001;
    }
  });

  return (
    <Sphere ref={meshRef} args={[2, 64, 64]}>
      <meshStandardMaterial
        color="#1e293b"
        roughness={0.7}
        metalness={0.2}
      />
    </Sphere>
  );
}

// Main Globe component
export function Globe() {
  const { cities, setSelectedCity } = useGlobeStore();

  const handleCityClick = (city: CityWithScore) => {
    setSelectedCity(city);
    console.log('Selected city:', city.name);
  };

  return (
    <div className="w-full h-screen">
      <Canvas
        camera={{ position: [0, 0, 5], fov: 45 }}
        gl={{ antialias: true }}
      >
        <ambientLight intensity={0.5} />
        <pointLight position={[10, 10, 10]} intensity={1} />

        <Earth />

        {cities.map((city) => (
          <CityMarker
            key={city.id}
            city={city}
            onClick={handleCityClick}
          />
        ))}

        <OrbitControls
          enableZoom={true}
          enablePan={false}
          minDistance={3}
          maxDistance={8}
        />
      </Canvas>
    </div>
  );
}
