export interface Channel {
  id: string;
  name: string;
  description: string;
  category: string;
  country: string;
  language: string;
  logo: string;
  streamUrl: string;
  featured: boolean;
}

// API Response Types
export interface ApiLogo {
  url: string;
  width: number;
  height: number;
}

export interface ApiCategory {
  name_1: string;
  name_2: string;
  name_3: string;
}

export interface ApiCountry {
  name: string;
}

export interface ApiLanguage {
  name: string;
}

export interface ApiFeaturedChannel {
  channel: string;
  title: string;
  url: string;
  quality: string;
  logo: ApiLogo | null;
  category: ApiCategory | null;
  country: ApiCountry | null;
  language: ApiLanguage | null;
}

// Fetch categories from backend API
export const fetchCategories = async (): Promise<string[]> => {
  try {
    const response = await fetch('http://127.0.0.1:8090/api/v1/stream/categories');
    if (!response.ok) throw new Error('Failed to fetch categories');
    const data = await response.json();
    return ['All', ...data.categories];
  } catch (error) {
    console.error('Error fetching categories:', error);
    return ['All'];
  }
};

// Fetch countries from backend API
export const fetchCountries = async (): Promise<string[]> => {
  try {
    const response = await fetch('http://127.0.0.1:8090/api/v1/stream/countries');
    if (!response.ok) throw new Error('Failed to fetch countries');
    const data = await response.json();
    return ['All', ...data.countries];
  } catch (error) {
    console.error('Error fetching countries:', error);
    return ['All'];
  }
};

// Fetch languages from backend API
export const fetchLanguages = async (): Promise<string[]> => {
  try {
    const response = await fetch('http://127.0.0.1:8090/api/v1/stream/languages');
    if (!response.ok) throw new Error('Failed to fetch languages');
    const data = await response.json();
    return ['All', ...data.languages];
  } catch (error) {
    console.error('Error fetching languages:', error);
    return ['All'];
  }
};

// Mock data - In production, this would come from an API
export const mockChannels: Channel[] = [
  {
    id: "1",
    name: "Global News Network",
    description: "24/7 breaking news coverage from around the world",
    category: "News",
    country: "USA",
    language: "English",
    logo: "https://images.unsplash.com/photo-1504711434969-e33886168f5c?w=200&h=200&fit=crop",
    streamUrl: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/BigBuckBunny.mp4",
    featured: true
  },
  {
    id: "2",
    name: "Sports Central",
    description: "Live sports action, highlights, and analysis",
    category: "Sports",
    country: "USA",
    language: "English",
    logo: "https://images.unsplash.com/photo-1461896836934-ffe607ba8211?w=200&h=200&fit=crop",
    streamUrl: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ElephantsDream.mp4",
    featured: true
  },
  {
    id: "3",
    name: "Cinema Plus",
    description: "Premium movies and entertainment 24/7",
    category: "Movies",
    country: "USA",
    language: "English",
    logo: "https://images.unsplash.com/photo-1489599849927-2ee91cede3ba?w=200&h=200&fit=crop",
    streamUrl: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerBlazes.mp4",
    featured: true
  },
  {
    id: "4",
    name: "Music TV",
    description: "Non-stop music videos and concerts",
    category: "Music",
    country: "USA",
    language: "English",
    logo: "https://images.unsplash.com/photo-1511379938547-c1f69419868d?w=200&h=200&fit=crop",
    streamUrl: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerEscapes.mp4",
    featured: true
  },
  {
    id: "5",
    name: "Kids World",
    description: "Safe and fun content for children",
    category: "Kids",
    country: "USA",
    language: "English",
    logo: "https://images.unsplash.com/photo-1503454537195-1dcabb73ffb9?w=200&h=200&fit=crop",
    streamUrl: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerFun.mp4",
    featured: true
  },
  {
    id: "6",
    name: "Discovery World",
    description: "Fascinating documentaries about our planet",
    category: "Documentary",
    country: "USA",
    language: "English",
    logo: "https://images.unsplash.com/photo-1451187580459-43490279c0fa?w=200&h=200&fit=crop",
    streamUrl: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerJoyrides.mp4",
    featured: true
  },
  {
    id: "7",
    name: "Tech Today",
    description: "Latest technology news and reviews",
    category: "Tech",
    country: "USA",
    language: "English",
    logo: "https://images.unsplash.com/photo-1518770660439-4636190af475?w=200&h=200&fit=crop",
    streamUrl: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerMeltdowns.mp4",
    featured: false
  },
  {
    id: "8",
    name: "Lifestyle Plus",
    description: "Cooking, fashion, and lifestyle shows",
    category: "Lifestyle",
    country: "USA",
    language: "English",
    logo: "https://images.unsplash.com/photo-1492562080023-ab3db95bfbce?w=200&h=200&fit=crop",
    streamUrl: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/Sintel.mp4",
    featured: false
  }
];

// Generate more mock channels to simulate 10,000+ channels
export const generateMockChannels = (): Channel[] => {
  const channels = [...mockChannels];
  const baseChannels = [...mockChannels];
  
  for (let i = 9; i <= 100; i++) {
    const baseChannel = baseChannels[i % baseChannels.length];
    channels.push({
      ...baseChannel,
      id: i.toString(),
      name: `${baseChannel.name} ${i}`,
      featured: false
    });
  }
  
  return channels;
};

export const allChannels = generateMockChannels();

export const getFeaturedChannels = () => allChannels.filter(c => c.featured);

export const getChannelsByCategory = (category: string) => 
  category === "All" ? allChannels : allChannels.filter(c => c.category === category);

export const getChannelsByCountry = (country: string) => 
  country === "All" ? allChannels : allChannels.filter(c => c.country === country);

export const getChannelsByLanguage = (language: string) => 
  language === "All" ? allChannels : allChannels.filter(c => c.language === language);

export const filterChannels = (category: string, country: string, language: string) => {
  let filtered = allChannels;
  if (category !== "All") filtered = filtered.filter(c => c.category === category);
  if (country !== "All") filtered = filtered.filter(c => c.country === country);
  if (language !== "All") filtered = filtered.filter(c => c.language === language);
  return filtered;
};

// Fetch featured channels from API
export const fetchFeaturedChannels = async (): Promise<Channel[]> => {
  try {
    const response = await fetch('http://127.0.0.1:8090/api/v1/stream/featured');
    if (!response.ok) {
      throw new Error('Failed to fetch featured channels');
    }
    const data: ApiFeaturedChannel[] = await response.json();
    
    // Convert API response to Channel interface
    return data.map((apiChannel, index) => ({
      id: apiChannel.channel, // Use channel name as ID for now
      name: apiChannel.title,
      description: apiChannel.title,
      category: apiChannel.category?.name_1 || 'General',
      country: apiChannel.country?.name || 'Unknown',
      language: apiChannel.language?.name || 'Unknown',
      logo: apiChannel.logo?.url || 'https://images.unsplash.com/photo-1504711434969-e33886168f5c?w=200&h=200&fit=crop',
      streamUrl: apiChannel.url,
      featured: true
    }));
  } catch (error) {
    console.error('Error fetching featured channels:', error);
    // Return empty array on error
    return [];
  }
};

// Fetch single channel data from API by channel name
export const fetchChannelByName = async (channelName: string): Promise<ApiFeaturedChannel | null> => {
  try {
    const response = await fetch(`http://127.0.0.1:8090/api/v1/stream/channel/${encodeURIComponent(channelName)}`);
    if (!response.ok) {
      if (response.status === 404) {
        return null;
      }
      throw new Error('Failed to fetch channel');
    }
    const data: ApiFeaturedChannel = await response.json();
    return data;
  } catch (error) {
    console.error('Error fetching channel:', error);
    return null;
  }
};

// Fetch channels by category
export const fetchChannelsByCategory = async (categoryName: string): Promise<Channel[]> => {
  try {
    const response = await fetch('http://127.0.0.1:8090/api/v1/stream/category', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ category_name: categoryName.toLowerCase() }),
    });
    
    if (!response.ok) {
      throw new Error('Failed to fetch channels by category');
    }
    
    const data: ApiFeaturedChannel[] = await response.json();
    
    return data.map((apiChannel) => ({
      id: apiChannel.channel,
      name: apiChannel.title,
      description: apiChannel.title,
      category: apiChannel.category?.name_1 || 'General',
      country: apiChannel.country?.name || 'Unknown',
      language: apiChannel.language?.name || 'Unknown',
      logo: apiChannel.logo?.url || 'https://images.unsplash.com/photo-1504711434969-e33886168f5c?w=200&h=200&fit=crop',
      streamUrl: apiChannel.url,
      featured: false
    }));
  } catch (error) {
    console.error('Error fetching channels by category:', error);
    return [];
  }
};

// Fetch recommended channels based on current watching channel
export const fetchRecommendedChannels = async (
  channel: string,
  categoryName: string,
  countryName: string,
  languageName: string
): Promise<Channel[]> => {
  try {
    const response = await fetch('http://127.0.0.1:8090/api/v1/stream/recommend', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        channel,
        category_name: categoryName,
        country_name: countryName,
        language_name: languageName,
      }),
    });
    
    if (!response.ok) {
      throw new Error('Failed to fetch recommended channels');
    }
    
    const data: ApiFeaturedChannel[] = await response.json();
    
    return data.map((apiChannel) => ({
      id: apiChannel.channel,
      name: apiChannel.title,
      description: apiChannel.title,
      category: apiChannel.category?.name_1 || 'General',
      country: apiChannel.country?.name || 'Unknown',
      language: apiChannel.language?.name || 'Unknown',
      logo: apiChannel.logo?.url || 'https://images.unsplash.com/photo-1504711434969-e33886168f5c?w=200&h=200&fit=crop',
      streamUrl: apiChannel.url,
      featured: false
    }));
  } catch (error) {
    console.error('Error fetching recommended channels:', error);
    return [];
  }
};

// Response type for all streams with pagination
export interface AllStreamsResponse {
  channels: Channel[];
  total: number;
  page: number;
  per_page: number;
  total_pages: number;
}

// Response type for search results
export interface SearchResponse {
  channels: Channel[];
  total: number;
}

// Fetch all streams with filtering and pagination
export const fetchAllStreams = async (
  category: string,
  country: string,
  language: string,
  page: number = 1
): Promise<AllStreamsResponse> => {
  try {
    const response = await fetch('http://127.0.0.1:8090/api/v1/stream/all', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        category: category.toLowerCase(),
        country,
        language,
        page,
      }),
    });
    
    if (!response.ok) {
      throw new Error('Failed to fetch all streams');
    }
    
    const data = await response.json();
    
    // Convert API response channels to frontend Channel format
    const channels = data.channels.map((apiChannel: ApiFeaturedChannel) => ({
      id: apiChannel.channel,
      name: apiChannel.title,
      description: apiChannel.title,
      category: apiChannel.category?.name_1 || 'General',
      country: apiChannel.country?.name || 'Unknown',
      language: apiChannel.language?.name || 'Unknown',
      logo: apiChannel.logo?.url || 'https://images.unsplash.com/photo-1504711434969-e33886168f5c?w=200&h=200&fit=crop',
      streamUrl: apiChannel.url,
      featured: false
    }));
    
    return {
      channels,
      total: data.total,
      page: data.page,
      per_page: data.per_page,
      total_pages: data.total_pages,
    };
  } catch (error) {
    console.error('Error fetching all streams:', error);
    return {
      channels: [],
      total: 0,
      page: 1,
      per_page: 24,
      total_pages: 0,
    };
  }
};

// Search for channels by name or title
export const searchChannels = async (query: string): Promise<SearchResponse> => {
  try {
    if (!query || query.trim() === '') {
      return { channels: [], total: 0 };
    }

    const response = await fetch('http://127.0.0.1:8090/api/v1/stream/search', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ query: query.trim() }),
    });

    if (!response.ok) {
      throw new Error('Failed to search channels');
    }

    const data = await response.json();

    // Convert API response channels to frontend Channel format
    const channels = data.channels.map((apiChannel: ApiFeaturedChannel) => ({
      id: apiChannel.channel,
      name: apiChannel.title,
      description: apiChannel.title,
      category: apiChannel.category?.name_1 || 'General',
      country: apiChannel.country?.name || 'Unknown',
      language: apiChannel.language?.name || 'Unknown',
      logo: apiChannel.logo?.url || 'https://images.unsplash.com/photo-1504711434969-e33886168f5c?w=200&h=200&fit=crop',
      streamUrl: apiChannel.url,
      featured: false
    }));

    return {
      channels,
      total: data.total,
    };
  } catch (error) {
    console.error('Error searching channels:', error);
    return {
      channels: [],
      total: 0,
    };
  }
};
