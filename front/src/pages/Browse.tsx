import { useState, useEffect } from "react";
import { useSearchParams } from "react-router-dom";
import Header from "@/components/Header";
import ChannelCard from "@/components/ChannelCard";
import VerticalAd from "@/components/ads/VerticalAd";
import { fetchCategories, fetchCountries, fetchLanguages, fetchAllStreams, searchChannels, Channel } from "@/lib/channels";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useLanguage } from "@/contexts/LanguageContext";

const Browse = () => {
  const [searchParams, setSearchParams] = useSearchParams();
  const [selectedCategory, setSelectedCategory] = useState("All");
  const [selectedCountry, setSelectedCountry] = useState("All");
  const [selectedLanguage, setSelectedLanguage] = useState("All");
  const [currentPage, setCurrentPage] = useState(1);
  const [displayedChannels, setDisplayedChannels] = useState<Channel[]>([]);
  const [totalChannels, setTotalChannels] = useState(0);
  const [totalPages, setTotalPages] = useState(0);
  const [loading, setLoading] = useState(false);
  const [filtersLoading, setFiltersLoading] = useState(true);
  const [categories, setCategories] = useState<string[]>(['All']);
  const [countries, setCountries] = useState<string[]>(['All']);
  const [languages, setLanguages] = useState<string[]>(['All']);
  const { t } = useLanguage();

  // Fetch filter options from PocketBase on mount
  useEffect(() => {
    const loadFilters = async () => {
      setFiltersLoading(true);
      try {
        const [cats, ctrs, langs] = await Promise.all([
          fetchCategories(),
          fetchCountries(),
          fetchLanguages()
        ]);
        setCategories(cats);
        setCountries(ctrs);
        setLanguages(langs);
      } catch (error) {
        console.error('Error loading filters:', error);
        // Keep default 'All' values if fetch fails
      } finally {
        setFiltersLoading(false);
      }
    };
    loadFilters();
  }, []);

  const searchQuery = searchParams.get("search") || "";

  useEffect(() => {
    const loadChannels = async () => {
      setLoading(true);
      try {
        if (searchQuery) {
          // Handle search functionality
          const searchResult = await searchChannels(searchQuery);
          setDisplayedChannels(searchResult.channels);
          setTotalChannels(searchResult.total);
          setTotalPages(1); // Search results don't have pagination
        } else {
          // Load normal channels with filters
          const result = await fetchAllStreams(
            selectedCategory,
            selectedCountry,
            selectedLanguage,
            currentPage
          );
          
          setDisplayedChannels(result.channels);
          setTotalChannels(result.total);
          setTotalPages(result.total_pages);
        }
      } catch (error) {
        console.error('Error loading channels:', error);
        setDisplayedChannels([]);
        setTotalChannels(0);
        setTotalPages(0);
      } finally {
        setLoading(false);
      }
    };

    loadChannels();
  }, [selectedCategory, selectedCountry, selectedLanguage, currentPage, searchQuery]);

  const handleSearch = (query: string) => {
    setSearchParams(query ? { search: query } : {});
    setCurrentPage(1);
  };

  const handleCategoryChange = (category: string) => {
    setSelectedCategory(category);
    setSearchParams({});
    setCurrentPage(1);
  };

  const handleCountryChange = (country: string) => {
    setSelectedCountry(country);
    setSearchParams({});
    setCurrentPage(1);
  };

  const handleLanguageChange = (language: string) => {
    setSelectedLanguage(language);
    setSearchParams({});
    setCurrentPage(1);
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

      <div className="container mx-auto py-8 px-4">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-4xl font-bold mb-2 text-foreground">
            {searchQuery ? `${t("browse.searchResults")} "${searchQuery}"` : t("browse.title")}
          </h1>
          <p className="text-muted-foreground">
            {t("browse.showing")} {displayedChannels.length} {t("browse.of")} {totalChannels.toLocaleString()} {t("browse.channels")}
          </p>
        </div>

        {/* Filters - Sticky */}
        <div className="sticky top-16 z-40 bg-background py-4 mb-8 flex flex-wrap gap-4 items-center border-b border-border">
          <div className="flex items-center gap-2">
            <span className="text-sm font-medium text-foreground">{t("browse.category")}:</span>
            <Select value={selectedCategory} onValueChange={handleCategoryChange} disabled={filtersLoading}>
              <SelectTrigger className="w-[180px] bg-secondary border-border">
                <SelectValue>
                  {filtersLoading ? "Loading..." : selectedCategory === "All" ? t("filter.all") : selectedCategory}
                </SelectValue>
              </SelectTrigger>
              <SelectContent className="bg-popover z-50">
                {categories.map((category) => (
                  <SelectItem key={category} value={category}>
                    {category === "All" ? t("filter.all") : category}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="flex items-center gap-2">
            <span className="text-sm font-medium text-foreground">{t("browse.country")}:</span>
            <Select value={selectedCountry} onValueChange={handleCountryChange} disabled={filtersLoading}>
              <SelectTrigger className="w-[180px] bg-secondary border-border">
                <SelectValue>
                  {filtersLoading ? "Loading..." : selectedCountry === "All" ? t("filter.all") : selectedCountry}
                </SelectValue>
              </SelectTrigger>
              <SelectContent className="bg-popover z-50">
                {countries.map((country) => (
                  <SelectItem key={country} value={country}>
                    {country === "All" ? t("filter.all") : country}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="flex items-center gap-2">
            <span className="text-sm font-medium text-foreground">{t("browse.language")}:</span>
            <Select value={selectedLanguage} onValueChange={handleLanguageChange} disabled={filtersLoading}>
              <SelectTrigger className="w-[180px] bg-secondary border-border">
                <SelectValue>
                  {filtersLoading ? "Loading..." : selectedLanguage === "All" ? t("filter.all") : selectedLanguage}
                </SelectValue>
              </SelectTrigger>
              <SelectContent className="bg-popover z-50">
                {languages.map((language) => (
                  <SelectItem key={language} value={language}>
                    {language === "All" ? t("filter.all") : language}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </div>

        {/* Channels Grid */}
        {loading ? (
          <div className="text-center py-16">
            <p className="text-xl text-muted-foreground">Loading channels...</p>
          </div>
        ) : displayedChannels.length > 0 ? (
          <>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6 mb-8">
              {displayedChannels.map((channel) => (
                <ChannelCard key={channel.id} channel={channel} />
              ))}
            </div>

            {/* Pagination */}
            {totalPages > 1 && (
              <div className="flex justify-center gap-2">
                <Button
                  variant="outline"
                  onClick={() => setCurrentPage((p) => Math.max(1, p - 1))}
                  disabled={currentPage === 1}
                >
                  {t("browse.previous")}
                </Button>
                <div className="flex items-center gap-2 px-4">
                  <span className="text-sm text-muted-foreground">
                    {t("browse.page")} {currentPage} {t("browse.of")} {totalPages}
                  </span>
                </div>
                <Button
                  variant="outline"
                  onClick={() => setCurrentPage((p) => Math.min(totalPages, p + 1))}
                  disabled={currentPage === totalPages}
                >
                  {t("browse.next")}
                </Button>
              </div>
            )}
          </>
        ) : (
          <div className="text-center py-16">
            <p className="text-xl text-muted-foreground">{t("browse.noChannels")}</p>
          </div>
        )}
      </div>
      </div>
    </div>
  );
};

export default Browse;
