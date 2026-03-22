import { useCallback, useMemo } from 'react';
import Map, { Marker, NavigationControl, AttributionControl } from 'react-map-gl/maplibre';
import 'maplibre-gl/dist/maplibre-gl.css';
import { useGlobeStore } from '../store';
import type { CityWithScore } from '../types';

// Esri World Imagery (free, no API key) + Esri reference labels overlay
// eslint-disable-next-line @typescript-eslint/no-explicit-any
const MAP_STYLE: any = {
  version: 8,
  glyphs: 'https://demotiles.maplibre.org/font/{fontstack}/{range}.pbf',
  fog: {
    color: 'rgb(10, 15, 30)',
    'high-color': 'rgb(30, 50, 100)',
    'horizon-blend': 0.03,
    'space-color': 'rgb(5, 8, 20)',
    'star-intensity': 0.6,
  },
  sources: {
    satellite: {
      type: 'raster',
      tiles: [
        'https://server.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile/{z}/{y}/{x}',
      ],
      tileSize: 256,
      maxzoom: 19,
      attribution: '© Esri, Maxar, Earthstar Geographics',
    },
    labels: {
      type: 'raster',
      tiles: [
        'https://services.arcgisonline.com/ArcGIS/rest/services/Reference/World_Boundaries_and_Places/MapServer/tile/{z}/{y}/{x}',
      ],
      tileSize: 256,
      maxzoom: 19,
    },
  },
  layers: [
    { id: 'satellite', type: 'raster', source: 'satellite' },
    { id: 'labels', type: 'raster', source: 'labels', minzoom: 1 },
  ],
};

function CityPin({ city, onClick }: { city: CityWithScore; onClick: (city: CityWithScore) => void }) {
  const color = useMemo(() => {
    if (city.score >= 85) return '#10b981';
    if (city.score >= 70) return '#f59e0b';
    return '#ef4444';
  }, [city.score]);

  return (
    <div
      onClick={() => onClick(city)}
      title={city.name}
      style={{
        width: 14,
        height: 14,
        borderRadius: '50%',
        backgroundColor: color,
        border: '2px solid white',
        boxShadow: `0 0 6px ${color}`,
        cursor: 'pointer',
        transform: 'translate(-50%, -50%)',
      }}
    />
  );
}

export function Globe() {
  const { cities, setSelectedCity } = useGlobeStore();

  const handleCityClick = useCallback((city: CityWithScore) => {
    setSelectedCity(city);
  }, [setSelectedCity]);

  return (
    <div className="w-full h-screen">
      <Map
        initialViewState={{ longitude: 20, latitude: 20, zoom: 1.5 }}
        style={{ width: '100%', height: '100%' }}
        mapStyle={MAP_STYLE}
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        projection={{ name: 'globe' } as any}
        attributionControl={false}
      >
        <NavigationControl position="bottom-right" showCompass={false} />
        <AttributionControl compact position="bottom-left" />

        {cities.map((city) => (
          <Marker
            key={city.id}
            longitude={city.lng}
            latitude={city.lat}
            anchor="center"
          >
            <CityPin city={city} onClick={handleCityClick} />
          </Marker>
        ))}
      </Map>
    </div>
  );
}
