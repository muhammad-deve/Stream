import { useParams, useNavigate, useLocation } from "react-router-dom";
import { useEffect, useState, useRef } from "react";
import Header from "@/components/Header";
import { fetchRecommendedChannels, fetchChannelByName, Channel, ApiFeaturedChannel, playStream } from "@/lib/channels";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { ArrowLeft, Globe, Languages } from "lucide-react";
import { useLanguage } from "@/contexts/LanguageContext";
import ChannelCard from "@/components/ChannelCard";
import Hls from "hls.js";

const ChannelPage = () => {
  const { id } = useParams();
  const location = useLocation();
  const navigate = useNavigate();
  const { t } = useLanguage();
  const [channel, setChannel] = useState<Channel | null>(null);
  const [recommendedChannels, setRecommendedChannels] = useState<Channel[]>([]);
  const [loading, setLoading] = useState(true);
  const videoRef = useRef<HTMLVideoElement>(null);
  const hlsRef = useRef<Hls | null>(null);

  useEffect(() => {
    const loadChannel = async () => {
      if (!id) return;
      setLoading(true);
      
      // Get channel data from navigation state
      let channelData = location.state?.channelData as Channel | undefined;
      
      // If no state, fetch from API
      if (!channelData) {
        const apiChannel = await fetchChannelByName(decodeURIComponent(id));
        if (!apiChannel) {
          // Channel not found, redirect to home
          navigate('/');
          return;
        }
        
        // Convert API channel to Channel format
        channelData = {
          id: apiChannel.channel,
          name: apiChannel.title,
          description: apiChannel.title,
          category: apiChannel.category?.name_1 || 'General',
          country: apiChannel.country?.name || 'Unknown',
          language: apiChannel.language?.name || 'Unknown',
          logo: apiChannel.logo?.url || 'https://images.unsplash.com/photo-1504711434969-e33886168f5c?w=200&h=200&fit=crop',
          streamUrl: apiChannel.url,
          featured: false
        };
      }
      
      setChannel(channelData);
      setLoading(false);
      
      // Fetch recommended channels
      const recommended = await fetchRecommendedChannels(
        channelData.id, // Use the channel name/id
        channelData.category,
        channelData.country,
        channelData.language
      );
      setRecommendedChannels(recommended);
    };
    loadChannel();
  }, [id, location.state]);

  // Setup HLS player
  useEffect(() => {
    if (!channel || !videoRef.current) return;

    const video = videoRef.current;
    const token = channel.streamUrl; // This is the UUID token

    // Cleanup previous instance
    if (hlsRef.current) {
      hlsRef.current.destroy();
      hlsRef.current = null;
    }

    // Resolve token to actual URL and play
    const setupPlayer = async () => {
      try {
        // Call backend to resolve token to actual URL
        const actualUrl = await playStream(token);
        
        if (!actualUrl) {
          console.error("Failed to resolve stream URL - token may be expired");
          // You could show an error message to user here
          return;
        }

        // Check if HLS is supported
        if (Hls.isSupported()) {
          const hls = new Hls({
            enableWorker: true,
            lowLatencyMode: true,
          });
          
          hlsRef.current = hls;
          hls.loadSource(actualUrl);
          hls.attachMedia(video);
          
          hls.on(Hls.Events.MANIFEST_PARSED, () => {
            video.play().catch(err => console.log("Autoplay prevented:", err));
          });

          hls.on(Hls.Events.ERROR, (event, data) => {
            if (data.fatal) {
              switch (data.type) {
                case Hls.ErrorTypes.NETWORK_ERROR:
                  console.error("Network error encountered, trying to recover...");
                  hls.startLoad();
                  break;
                case Hls.ErrorTypes.MEDIA_ERROR:
                  console.error("Media error encountered, trying to recover...");
                  hls.recoverMediaError();
                  break;
                default:
                  console.error("Fatal error, cannot recover");
                  hls.destroy();
                  break;
              }
            }
          });
        } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
          // For Safari (native HLS support)
          video.src = actualUrl;
          video.addEventListener('loadedmetadata', () => {
            video.play().catch(err => console.log("Autoplay prevented:", err));
          });
        }
      } catch (error) {
        console.error("Error setting up player:", error);
      }
    };

    setupPlayer();

    // Cleanup on unmount
    return () => {
      if (hlsRef.current) {
        hlsRef.current.destroy();
        hlsRef.current = null;
      }
    };
  }, [channel]);

  if (loading) {
    return (
      <div className="min-h-screen bg-background">
        <Header onSearch={(q) => navigate(`/browse?search=${encodeURIComponent(q)}`)} />
        <div className="container mx-auto py-16 px-4 text-center">
          <p className="text-muted-foreground">Loading channel...</p>
        </div>
      </div>
    );
  }

  if (!channel) {
    return (
      <div className="min-h-screen bg-background">
        <Header onSearch={(q) => navigate(`/browse?search=${encodeURIComponent(q)}`)} />
        <div className="container mx-auto py-16 px-4 text-center">
          <h1 className="text-4xl font-bold mb-4 text-foreground">{t("channel.notFound")}</h1>
          <Button onClick={() => navigate("/")}>{t("channel.goHome")}</Button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background">
      <Header onSearch={(q) => navigate(`/browse?search=${encodeURIComponent(q)}`)} />

      <div className="container mx-auto max-w-7xl py-8 px-4">
        {/* Back Button */}
        <Button
          variant="ghost"
          className="mb-6"
          onClick={() => navigate(-1)}
        >
          <ArrowLeft className="h-4 w-4 mr-2" />
          {t("channel.back")}
        </Button>

        {/* Video Player */}
        <Card className="overflow-hidden bg-card border-border mb-6">
          <div className="aspect-video bg-black">
            <video
              ref={videoRef}
              controls
              className="w-full h-full"
              poster={channel.logo}
            >
              Your browser does not support the video tag.
            </video>
          </div>
        </Card>

        {/* Channel Info */}
        <div className="mb-8">
          <div className="flex items-start justify-between gap-4 mb-4">
            <div className="flex-1">
              <h1 className="text-3xl font-bold mb-2 text-foreground">{channel.name}</h1>
              <div className="flex flex-wrap items-center gap-3 text-sm text-muted-foreground mb-3">
                <div className="flex items-center gap-1">
                  <Globe className="h-4 w-4" />
                  <span>{channel.country}</span>
                </div>
                <div className="flex items-center gap-1">
                  <Languages className="h-4 w-4" />
                  <span>{channel.language}</span>
                </div>
                <span className="px-3 py-1 rounded-full bg-primary/20 text-primary">
                  {channel.category}
                </span>
              </div>
            </div>
          </div>

          <p className="text-muted-foreground leading-relaxed">{channel.description}</p>
        </div>

        {/* Recommended Channels */}
        {recommendedChannels.length > 0 && (
          <div className="mt-12">
            <h2 className="text-2xl font-bold mb-6 text-foreground">
              Recommended Channels
            </h2>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
              {recommendedChannels.map((c) => (
                <ChannelCard key={c.id} channel={c} />
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default ChannelPage;
