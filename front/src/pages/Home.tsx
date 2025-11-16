import { useState, useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";
import Header from "@/components/Header";
import ChannelCard from "@/components/ChannelCard";
import VerticalAd from "@/components/ads/VerticalAd";
import HorizontalAd from "@/components/ads/HorizontalAd";
import { fetchFeaturedChannels, fetchChannelsByCategory, fetchCategories, Channel } from "@/lib/channels";
import { Button } from "@/components/ui/button";
import { Sparkles } from "lucide-react";
import { useLanguage } from "@/contexts/LanguageContext";

const Home = () => {
  const [selectedCategory, setSelectedCategory] = useState("All");
  const [featuredChannels, setFeaturedChannels] = useState<Channel[]>([]);
  const [categoryChannels, setCategoryChannels] = useState<Channel[]>([]);
  const [categories, setCategories] = useState<string[]>(["All"]);
  const [loading, setLoading] = useState(true);
  const [categoryLoading, setCategoryLoading] = useState(false);
  const [mousePos, setMousePos] = useState({ x: 0, y: 0 });
  const heroRef = useRef<HTMLDivElement>(null);
  const navigate = useNavigate();
  const { t } = useLanguage();

  // Track mouse movement for interactive effects
  useEffect(() => {
    const handleMouseMove = (e: MouseEvent) => {
      if (heroRef.current) {
        const rect = heroRef.current.getBoundingClientRect();
        const x = e.clientX - rect.left;
        const y = e.clientY - rect.top;
        setMousePos({ x, y });
      }
    };

    window.addEventListener("mousemove", handleMouseMove);
    return () => window.removeEventListener("mousemove", handleMouseMove);
  }, []);

  // Fetch categories and featured channels from API on component mount
  useEffect(() => {
    const loadData = async () => {
      setLoading(true);
      try {
        const [cats, channels] = await Promise.all([
          fetchCategories(),
          fetchFeaturedChannels()
        ]);
        setCategories(cats);
        setFeaturedChannels(channels);
      } catch (error) {
        console.error('Error loading data:', error);
        // Keep default values if fetch fails
        setCategories(["All"]);
        setFeaturedChannels([]);
      } finally {
        setLoading(false);
      }
    };
    loadData();
  }, []);

  // Fetch channels by category when category changes
  useEffect(() => {
    const loadCategoryChannels = async () => {
      setCategoryLoading(true);
      try {
        const channels = await fetchChannelsByCategory(selectedCategory);
        setCategoryChannels(channels);
      } catch (error) {
        console.error('Error loading category channels:', error);
        setCategoryChannels([]);
      } finally {
        setCategoryLoading(false);
      }
    };
    loadCategoryChannels();
  }, [selectedCategory]);

  const handleSearch = (query: string) => {
    if (query.trim()) {
      navigate(`/browse?search=${encodeURIComponent(query)}`);
    }
  };

  return (
    <div className="min-h-screen bg-background relative">
      {/* Left Vertical Ad */}
      <div className="fixed left-0 top-20 hidden xl:block z-10">
        <VerticalAd />
      </div>
      
      {/* Right Vertical Ad */}
      <div className="fixed right-0 top-20 hidden xl:block z-10">
        <VerticalAd />
      </div>

      {/* Main Content with side margins */}
      <div className="xl:mx-[180px]">
        <Header onSearch={handleSearch} />

      {/* Hero Section */}
      <section ref={heroRef} className="relative overflow-hidden py-32 px-4 group">
        {/* Background */}
        <div className="absolute inset-0 bg-background"></div>
        
        {/* Streaming Device Icons - Left Side */}
        <div className="absolute left-0 top-1/4 opacity-20 dark:opacity-30 group-hover:opacity-40 transition-all duration-500 pointer-events-none transform group-hover:scale-110 -translate-x-1/3 sm:-translate-x-1/4 md:translate-x-0">
          <svg className="w-32 sm:w-48 md:w-64 h-32 sm:h-48 md:h-64 text-blue-500 dark:text-blue-400" fill="currentColor" viewBox="0 0 24 24">
            {/* TV Icon */}
            <path d="M20 3H4c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14l4 4V5c0-1.1-.9-2-2-2zm-2 16H4V5h14v14z"/>
          </svg>
        </div>
        
        {/* Play Button Icon - Right Side */}
        <div className="absolute right-0 bottom-1/4 opacity-20 dark:opacity-30 group-hover:opacity-40 transition-all duration-500 pointer-events-none transform group-hover:scale-110 translate-x-1/3 sm:translate-x-1/4 md:translate-x-0">
          <svg className="w-40 sm:w-56 md:w-80 h-40 sm:h-56 md:h-80 text-purple-500 dark:text-purple-400" fill="currentColor" viewBox="0 0 24 24">
            {/* Play Icon */}
            <path d="M8 5v14l11-7z"/>
          </svg>
        </div>
        
        {/* Streaming Waves - Center */}
        <div className="absolute inset-0 flex items-center justify-center opacity-10 dark:opacity-20 group-hover:opacity-25 transition-all duration-500 pointer-events-none transform group-hover:scale-105">
          <svg className="w-48 sm:w-64 md:w-96 h-48 sm:h-64 md:h-96 text-cyan-500 dark:text-cyan-400" viewBox="0 0 100 100" fill="none" stroke="currentColor" strokeWidth="2">
            <circle cx="50" cy="50" r="10"/>
            <circle cx="50" cy="50" r="25" opacity="0.7"/>
            <circle cx="50" cy="50" r="40" opacity="0.4"/>
            <circle cx="50" cy="50" r="55" opacity="0.2"/>
          </svg>
        </div>
        
        {/* Interactive Hover Gradient Overlay */}
        <div 
          className="absolute inset-0 opacity-0 group-hover:opacity-10 transition-all duration-500 pointer-events-none"
          style={{
            background: `radial-gradient(circle at ${mousePos.x}px ${mousePos.y}px, rgba(59, 130, 246, 0.4) 0%, transparent 50%)`
          }}
        ></div>
        
        <div className="container mx-auto text-center relative z-10">
          <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-primary/10 border border-primary/20 mb-6">
            <Sparkles className="h-4 w-4 text-primary" />
            <span className="text-sm font-medium text-primary">{t("hero.badge")}</span>
          </div>
          <h1 className="text-4xl md:text-6xl font-bold mb-4 text-foreground">
            {t("hero.title")}
            <br />
            <span className="text-primary">{t("hero.titleHighlight")}</span>
          </h1>
          <p className="text-xl text-muted-foreground mb-8 max-w-2xl mx-auto">
            {t("hero.description")}
          </p>
          <div className="flex flex-wrap gap-4 justify-center">
            <Button
              size="lg"
              className="bg-primary hover:bg-primary/90 text-primary-foreground"
              onClick={() => navigate("/browse")}
            >
              {t("hero.browseButton")}
            </Button>
            <Button
              size="lg"
              variant="outline"
              className="border-border hover:bg-secondary"
              onClick={() => navigate("/donate")}
            >
              {t("hero.supportButton")}
            </Button>
          </div>
        </div>
      </section>

      {/* Featured Channels */}
      <section className="py-16 px-4">
        <div className="container mx-auto max-w-7xl px-6">
          <h2 className="text-3xl font-bold mb-8 text-foreground">{t("section.featured")}</h2>
          {loading ? (
            <div className="text-center py-12">
              <p className="text-muted-foreground">Loading featured channels...</p>
            </div>
          ) : featuredChannels.length === 0 ? (
            <div className="text-center py-12">
              <p className="text-muted-foreground">No featured channels available</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-8">
              {featuredChannels.map((channel) => (
                <ChannelCard key={channel.id} channel={channel} />
              ))}
            </div>
          )}
        </div>
      </section>

      {/* Browse by Category */}
      <section className="py-16 px-4 bg-secondary/30">
        <div className="container mx-auto">
          <h2 className="text-3xl font-bold mb-8 text-foreground">{t("section.browseCategory")}</h2>
          
          {/* Category Filters */}
          <div className="flex flex-wrap gap-2 mb-8">
            {categories.map((category) => (
              <Button
                key={category}
                variant={selectedCategory === category ? "default" : "outline"}
                onClick={() => setSelectedCategory(category)}
                className={selectedCategory === category ? "bg-primary text-primary-foreground" : ""}
              >
                {category}
              </Button>
            ))}
          </div>

          {/* Channels Grid */}
          {categoryLoading ? (
            <div className="text-center py-12">
              <p className="text-muted-foreground">Loading channels...</p>
            </div>
          ) : categoryChannels.length === 0 ? (
            <div className="text-center py-12">
              <p className="text-muted-foreground">No channels found in this category</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
              {categoryChannels.map((channel) => (
                <ChannelCard key={channel.id} channel={channel} />
              ))}
            </div>
          )}

          <div className="text-center mt-8">
            <Button
              size="lg"
              variant="outline"
              onClick={() => navigate("/browse")}
            >
              {t("section.viewAll")}
            </Button>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="bg-secondary py-12 px-4 mt-16">
        <div className="container mx-auto">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8 mb-8">
            <div>
              <h3 className="font-bold text-lg mb-4 text-foreground">{t("footer.about")}</h3>
              <p className="text-muted-foreground text-sm">
                {t("footer.aboutText")}
              </p>
            </div>
            <div>
              <h3 className="font-bold text-lg mb-4 text-foreground">{t("footer.quickLinks")}</h3>
              <ul className="space-y-2 text-sm text-muted-foreground">
                <li><a href="/" className="hover:text-primary transition-colors">{t("header.home")}</a></li>
                <li><a href="/browse" className="hover:text-primary transition-colors">{t("header.browse")}</a></li>
                <li><a href="/donate" className="hover:text-primary transition-colors">{t("header.donate")}</a></li>
              </ul>
            </div>
            <div>
              <h3 className="font-bold text-lg mb-4 text-foreground">{t("footer.support")}</h3>
              <p className="text-muted-foreground text-sm mb-4">
                {t("footer.supportText")}
              </p>
              <Button variant="outline" onClick={() => navigate("/donate")}>
                {t("footer.makeDonation")}
              </Button>
            </div>
          </div>
          <div className="border-t border-border pt-8 text-center text-sm text-muted-foreground space-y-2">
            <p>{t("footer.copyright")}</p>
            <p className="text-xs">
              {t("footer.attribution")} <a href="https://github.com/iptv-org" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">iptv-org</a>
            </p>
          </div>
        </div>
      </footer>
      </div>
      
      {/* Bottom Horizontal Ad */}
      <div className="fixed bottom-0 left-0 right-0 z-20 bg-background border-t border-border">
        <div className="xl:mx-[180px]">
          <HorizontalAd />
        </div>
      </div>
    </div>
  );
};

export default Home;
