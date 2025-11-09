import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import Header from "@/components/Header";
import ChannelCard from "@/components/ChannelCard";
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
  const navigate = useNavigate();
  const { t } = useLanguage();

  // Fetch categories and featured channels from API on component mount
  useEffect(() => {
    const loadData = async () => {
      setLoading(true);
      const [cats, channels] = await Promise.all([
        fetchCategories(),
        fetchFeaturedChannels()
      ]);
      setCategories(cats);
      setFeaturedChannels(channels);
      setLoading(false);
    };
    loadData();
  }, []);

  // Fetch channels by category when category changes
  useEffect(() => {
    const loadCategoryChannels = async () => {
      setCategoryLoading(true);
      const channels = await fetchChannelsByCategory(selectedCategory);
      setCategoryChannels(channels);
      setCategoryLoading(false);
    };
    loadCategoryChannels();
  }, [selectedCategory]);

  const handleSearch = (query: string) => {
    if (query.trim()) {
      navigate(`/browse?search=${encodeURIComponent(query)}`);
    }
  };

  return (
    <div className="min-h-screen bg-background">
      <Header onSearch={handleSearch} />

      {/* Hero Section */}
      <section className="relative overflow-hidden py-20 px-4">
        <div className="absolute inset-0 bg-gradient-to-br from-background via-secondary to-background opacity-90"></div>
        <div className="absolute inset-0 opacity-10" style={{
          backgroundImage: `url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23ffffff' fill-opacity='0.4'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E")`
        }}></div>
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
  );
};

export default Home;
