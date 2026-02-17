package location

import "log/slog"

type Location struct {
	City        string
	CountryCode string
	CountryName string
	TimeZone    string
	IsDesktop   *bool
	IsMobile    *bool
	IsTablet    *bool
}

func (l *Location) LogValue() slog.Value {
	var attrs []slog.Attr
	if l.City != "" {
		attrs = append(attrs, slog.String("city", l.City))
	}

	if l.CountryCode != "" {
		attrs = append(attrs, slog.String("country_code", l.CountryCode))
	}

	if l.CountryName != "" {
		attrs = append(attrs, slog.String("country_name", l.CountryName))
	}

	if l.TimeZone != "" {
		attrs = append(attrs, slog.String("time_zone", l.TimeZone))
	}

	if l.IsDesktop != nil {
		attrs = append(attrs, slog.Bool("is_desktop", *l.IsDesktop))
	}

	if l.IsMobile != nil {
		attrs = append(attrs, slog.Bool("is_mobile", *l.IsMobile))
	}

	if l.IsTablet != nil {
		attrs = append(attrs, slog.Bool("is_tablet", *l.IsTablet))
	}

	return slog.GroupValue(attrs...)
}

func (l *Location) String() string {
	return l.LogValue().String()
}
